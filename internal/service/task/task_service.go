package task

import (
	"context"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/types"
	"github.com/eininst/scheduler/internal/util"
	"github.com/go-redis/redis/v8"
	"github.com/jinzhu/copier"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

const (
	STATUS_STOP = "stop"
	STATUS_RUN  = "run"
	STATUS_DEL  = "del"
)

type TaskService interface {
	Add(ctx context.Context, task *model.SchedulerTask) error
	Update(ctx context.Context, task *model.SchedulerTask) error
	PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*types.TaskDTO], error)
	Start(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	Del(ctx context.Context, id int64) error
}
type taskService struct {
	parse cron.Parser
	DB    *gorm.DB      `inject:""`
	Rcli  *redis.Client `inject:""`
	RsCli rs.Client     `inject:""`
}

func NewService() TaskService {
	return &taskService{
		parse: cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		),
	}
}

func (t *taskService) Add(ctx context.Context, task *model.SchedulerTask) error {
	_, er := t.parse.Parse(task.Spec)
	if er != nil {
		return service.NewServiceError("无效的Cron表达式")
	}
	if task.Timeout == 0 {
		task.Timeout = 15
	}
	task.Status = STATUS_STOP
	task.CreateTime = util.FormatTime()

	session := t.DB.WithContext(ctx)

	var ckTask model.SchedulerTask

	session.First(&ckTask, "name = ? and status != ?", task.Name, STATUS_DEL)
	if ckTask.Id != 0 {
		return service.NewServiceError("任务名称已存在")
	}
	err := session.Create(&task).Error
	if err != nil {
		return service.NewServiceError("创建任务失败")
	}

	return nil
}

func (t *taskService) Update(ctx context.Context, task *model.SchedulerTask) error {
	session := t.DB.WithContext(ctx)

	var _task model.SchedulerTask
	session.First(&_task, task.Id)
	if _task.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	_, er := t.parse.Parse(task.Spec)
	if er != nil {
		return service.NewServiceError("无效的Cron表达式")
	}
	if task.Timeout == 0 {
		task.Timeout = 15
	}
	var ckTask model.SchedulerTask

	session.First(&ckTask, "name = ? and id != ?", task.Name, _task.Id)
	if ckTask.Id != 0 {
		return service.NewServiceError("任务名称已存在")
	}
	err := session.Model(&_task).Updates(&model.SchedulerTask{
		Name:        task.Name,
		Group:       task.Group,
		Spec:        task.Spec,
		Url:         task.Url,
		Method:      task.Method,
		ContentType: task.ContentType,
		Body:        task.Body,
		Timeout:     task.Timeout,
		MaxRetries:  task.MaxRetries,
		Desc:        task.Desc,
	}).Error
	if err != nil {
		return service.NewServiceError("修改任务失败")
	}

	return nil
}

func (ts *taskService) Start(ctx context.Context, id int64) error {
	er := ts.updateStatus(ctx, id, STATUS_RUN)
	if er != nil {
		return service.NewServiceError("启动失败")
	}

	er = ts.RsCli.Send("task_register", rs.H{
		"task_id": id,
	})
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) Stop(ctx context.Context, id int64) error {
	er := ts.updateStatus(ctx, id, STATUS_STOP)
	if er != nil {
		return service.NewServiceError("停止失败")
	}
	er = ts.RsCli.Send("task_stop", rs.H{
		"task_id": id,
	})
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) Del(ctx context.Context, id int64) error {
	var t model.SchedulerTask

	sess := ts.DB.WithContext(ctx)
	sess.First(&t, id)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	er := sess.Delete(&t).Error
	if er != nil {
		return service.NewServiceError("删除失败")
	}

	er = ts.RsCli.Send("task_stop", rs.H{
		"task_id": id,
	})
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) updateStatus(ctx context.Context, id int64, status string) error {
	var t model.SchedulerTask

	sess := ts.DB.WithContext(ctx)
	sess.First(&t, id)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	return sess.Model(&t).Update("status", status).Error
}

func (t *taskService) PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*types.TaskDTO], error) {
	var total int64
	var tasks []*model.SchedulerTask

	sess := t.DB.WithContext(ctx)

	offset := (opt.Current - 1) * opt.PageSize
	if offset < 0 {
		offset = 0
	}
	q := sess.Model(&model.SchedulerTask{})
	if opt.Name != "" {
		q = q.Where("name LIKE ?", "%"+opt.Name+"%")
	}

	if opt.Group != "" {
		q = q.Where("`group` LIKE ?", "%"+opt.Group+"%")
	}

	if opt.UserId != 0 {
		q = q.Where("user_id = ?", opt.UserId)
	}

	if opt.Status != "" {
		q = q.Where("status = ?", opt.Status)
	}

	q.Count(&total)
	q.Limit(opt.PageSize).Offset(offset).Order("id desc").Find(&tasks)

	taskDtos := []*types.TaskDTO{}
	err := copier.Copy(&taskDtos, &tasks)
	if err != nil {
		return nil, err
	}

	userIds := util.ConvertIds(tasks, "UserId")

	var userList []*model.SchedulerUser
	sess.Find(&userList, userIds)
	userMap := util.ConvertIdMap(userList)

	for _, t := range taskDtos {
		if v, ok := userMap[t.UserId]; ok {
			t.UserName = v.Name
			t.UserRealName = v.RealName
			t.UserHead = v.Head
			t.UserMail = v.Mail
		}
	}

	pg := &types.Page[*types.TaskDTO]{
		Total: total,
		List:  taskDtos,
	}

	return pg, nil
}
