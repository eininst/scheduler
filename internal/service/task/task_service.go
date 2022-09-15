package task

import (
	"context"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/types"
	"github.com/eininst/scheduler/internal/util"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type TaskService interface {
	Add(ctx context.Context, task *model.SchedulerTask) error
	PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*model.SchedulerTask], error)
}
type taskService struct {
	parse cron.Parser
	DB    *gorm.DB      `inject:""`
	Rcli  *redis.Client `inject:""`
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
	task.Status = "paused"
	task.CreateTime = util.FormatTime()

	session := t.DB.WithContext(ctx)

	var ckTask model.SchedulerTask
	session.First(&ckTask, "name = ?", task.Name)
	if ckTask.Id != 0 {
		return service.NewServiceError("任务名称已存在")
	}
	err := session.Create(&task).Error
	if err != nil {
		return service.NewServiceError("创建任务失败")
	}

	return nil
}

func (t *taskService) PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*model.SchedulerTask], error) {
	var total int64
	var tasks []*model.SchedulerTask

	sess := t.DB.WithContext(ctx)

	offset := (opt.Current - 1) * opt.PageSize
	if offset < 0 {
		offset = 0
	}
	q := sess.Model(&model.SchedulerTask{})
	if opt.Name != "" {

	}

	if opt.UserId != 0 {
		q = q.Where("user_id = ?", opt.UserId)
	}

	q.Count(&total)
	q.Limit(opt.PageSize).Offset(offset).Order("id desc").Find(&tasks)

	pg := &types.Page[*model.SchedulerTask]{
		Total: total,
		List:  tasks,
	}
	return pg, nil
}
