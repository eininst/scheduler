package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eininst/flog"
	"github.com/eininst/rlock"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/eininst/scheduler/internal/types"
	"github.com/eininst/scheduler/internal/util"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/jinzhu/copier"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"sync"
	"time"
)

const (
	STATUS_STOP = "stop"
	STATUS_RUN  = "run"
	STATUS_DEL  = "del"

	GET    = "get"
	POST   = "post"
	PUT    = "put"
	DELETE = "delete"
)

type TaskService interface {
	RunTask(ctx context.Context)

	Add(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, task *model.Task) error
	PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*types.TaskDTO], error)
	Start(ctx context.Context, id int64) error
	Stop(ctx context.Context, id int64) error
	Del(ctx context.Context, id int64) error

	UpdateUser(ctx context.Context, tcu *types.TaskChangeUser) (int64, error)
	StartBatch(ctx context.Context, tbatch *types.TaskBatch) (int64, error)
	StopBatch(ctx context.Context, tbatch *types.TaskBatch) (int64, error)
	DelBatch(ctx context.Context, tbatch *types.TaskBatch) (int64, error)

	AddExcute(ctx context.Context, taskExcute *model.TaskExcute) error

	texcutePageByOption(ctx context.Context, opt *types.TaskExcuteOption) (*types.Page[*model.TaskExcute], error)
}
type taskService struct {
	parse   cron.Parser
	mux     *sync.Mutex
	TaskMap map[int64]cron.EntryID
	DB      *gorm.DB      `inject:""`
	Rcli    *redis.Client `inject:""`
	CronCli *cron.Cron    `inject:""`
	Rlock   *rlock.Rlock  `inject:""`

	UserService user.UserService `inject:""`
}

func NewService() TaskService {
	return &taskService{
		parse: cron.NewParser(
			cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
		),
		TaskMap: make(map[int64]cron.EntryID),
		mux:     &sync.Mutex{},
	}
}

func (t *taskService) Add(ctx context.Context, task *model.Task) error {
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

	var ckTask model.Task

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

func (t *taskService) Update(ctx context.Context, task *model.Task) error {
	session := t.DB.WithContext(ctx)

	var _task model.Task
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
	var ckTask model.Task

	session.First(&ckTask, "name = ? and id != ?", task.Name, _task.Id)
	if ckTask.Id != 0 {
		return service.NewServiceError("任务名称已存在")
	}
	err := session.Model(&_task).Updates(&model.Task{
		UserId:      task.UserId,
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

func (t *taskService) UpdateUser(ctx context.Context, tcu *types.TaskChangeUser) (int64, error) {
	if len(tcu.TaskIds) == 0 {
		return 0, nil
	}
	tk := &model.Task{}
	sqlp := fmt.Sprintf("update `%s` set user_id = ? where id in ?", tk.TableName())
	count := t.DB.WithContext(ctx).Exec(sqlp, tcu.UserId, tcu.TaskIds).RowsAffected

	return count, nil
}

func (t *taskService) StartBatch(ctx context.Context, tbatch *types.TaskBatch) (int64, error) {
	if len(tbatch.TaskIds) == 0 {
		return 0, service.NewServiceError("请选择任务")
	}

	count := int64(0)
	for _, id := range tbatch.TaskIds {
		er := t.Start(ctx, id)
		if er != nil {
			continue
		}
		count += 1
	}
	return count, nil
}

func (t *taskService) StopBatch(ctx context.Context, tbatch *types.TaskBatch) (int64, error) {
	if len(tbatch.TaskIds) == 0 {
		return 0, service.NewServiceError("请选择任务")
	}

	count := int64(0)
	for _, id := range tbatch.TaskIds {
		er := t.Stop(ctx, id)
		if er != nil {
			continue
		}
		count += 1
	}
	return count, nil
}

func (t *taskService) DelBatch(ctx context.Context, tbatch *types.TaskBatch) (int64, error) {
	if len(tbatch.TaskIds) == 0 {
		return 0, service.NewServiceError("请选择任务")
	}

	count := int64(0)
	for _, id := range tbatch.TaskIds {
		er := t.Del(ctx, id)
		if er != nil {
			continue
		}
		count += 1
	}
	return count, nil
}

func (ts *taskService) Start(ctx context.Context, id int64) error {
	ok, cancel := ts.Rlock.TryAcquire(fmt.Sprintf("task:%v", id), time.Second*10)
	defer cancel()
	if !ok {
		return service.NewServiceError("启动失败, 当前操作有其他人正在进行操作")
	}

	var t model.Task

	sess := ts.DB.WithContext(ctx)
	sess.First(&t, id)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}
	if t.Status == STATUS_RUN {
		return service.NewServiceError("任务已启动，请勿重复操作")
	}
	er := sess.Model(&t).Update("status", STATUS_RUN).Error
	if er != nil {
		return service.NewServiceError("启动失败")
	}

	er = ts.runTask(ctx, &t)
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) Stop(ctx context.Context, id int64) error {
	ok, cancel := ts.Rlock.TryAcquire(fmt.Sprintf("task:%v", id), time.Second*10)
	defer cancel()
	if !ok {
		return service.NewServiceError("停止失败, 当前任务有其他人正在进行操作")
	}

	var t model.Task

	sess := ts.DB.WithContext(ctx)
	sess.First(&t, id)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}
	if t.Status == STATUS_STOP {
		return service.NewServiceError("任务已停止，请勿重复操作")
	}

	er := sess.Model(&t).Update("status", STATUS_STOP).Error
	if er != nil {
		return service.NewServiceError("停止失败")
	}
	ts.delEntry(ctx, id)

	return nil
}

func (ts *taskService) Del(ctx context.Context, id int64) error {
	ok, cancel := ts.Rlock.TryAcquire(fmt.Sprintf("task:%v", id), time.Second*10)
	defer cancel()
	if !ok {
		return service.NewServiceError("删除失败, 当前任务有其他人正在进行操作")
	}

	var t model.Task

	sess := ts.DB.WithContext(ctx)
	sess.First(&t, id)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	er := sess.Delete(&t).Error
	if er != nil {
		return service.NewServiceError("删除失败")
	}

	ts.delEntry(ctx, id)

	return nil
}

func (ts *taskService) runTask(ctx context.Context, t *model.Task) error {
	flog.Info(t.Spec)
	eid, er := ts.CronCli.AddFunc(t.Spec, func() {
		ts.do(ctx, t)
	})
	flog.Info(eid)
	if er != nil {
		return er
	}

	ts.addEntry(ctx, t.Id, eid)
	return nil
}

func (t *taskService) PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*types.TaskDTO], error) {
	var total int64
	var tasks []*model.Task

	sess := t.DB.WithContext(ctx)

	offset := (opt.Current - 1) * opt.PageSize
	if offset < 0 {
		offset = 0
	}
	q := sess.Model(&model.Task{})
	if opt.Name != "" {
		q = q.Where("name LIKE ?", "%"+opt.Name+"%")
	}

	if opt.Group != "" {
		q = q.Where("`group` LIKE ?", "%"+opt.Group+"%")
	}

	if opt.UserId > 0 {
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

	var userList []*model.User
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

func (t *taskService) AddExcute(ctx context.Context, taskExcute *model.TaskExcute) error {
	session := t.DB.WithContext(ctx)
	taskExcute.CreateTime = util.FormatTime()

	err := session.Create(&taskExcute).Error
	if err != nil {
		return service.NewServiceError("创建任务记录失败")
	}

	return nil
}

func (ts *taskService) addEntry(ctx context.Context, taskId int64, eid cron.EntryID) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	ts.TaskMap[taskId] = eid
}

func (ts *taskService) delEntry(ctx context.Context, taskId int64) {
	ts.mux.Lock()
	defer ts.mux.Unlock()

	if eid, ok := ts.TaskMap[taskId]; ok {
		ts.CronCli.Remove(eid)
		delete(ts.TaskMap, taskId)
	}
}

func (ts *taskService) RunTask(ctx context.Context) {
	var tasks []*model.Task
	ts.DB.WithContext(ctx).Where("status = ?", STATUS_RUN).Find(&tasks)

	for _, t := range tasks {
		er := ts.runTask(ctx, t)
		flog.Info(er)
		if er != nil {
			flog.Fatal(er)
		}
	}
}

func (ts *taskService) TasktAlarm(ctx context.Context, mail, body string) {
	flog.Info(mail, body)
}

func (ts *taskService) do(ctx context.Context, t *model.Task) {
	jsonInfo, _ := json.Marshal(t)
	texcute := &model.TaskExcute{
		UserId:    t.UserId,
		TaskId:    t.Id,
		TaskName:  t.Name,
		TaskGroup: t.Group,
		TaskUrl:   t.Url,
		TaskObj:   string(jsonInfo),
		StartTime: util.FormatTimeMill(),
	}

	startDuration := time.Now().UnixMilli()

	cli := fiber.AcquireClient()

	var agt *fiber.Agent
	switch t.Method {
	case GET:
		agt = cli.Get(t.Url)
	case POST:
		agt = cli.Post(t.Url)
	case PUT:
		agt = cli.Put(t.Url)
	case DELETE:
		agt = cli.Delete(t.Url)
	default:
		flog.Error("无效的Method")
		return
	}

	agt.BodyString(t.Body)
	agt.Timeout(time.Second * time.Duration(t.Timeout))
	//agt.ReadTimeout = time.Second * 3
	agt.HostClient.MaxIdemponentCallAttempts = t.MaxRetries

	code, body, errors := agt.String()
	if len(errors) > 0 {
		elist := []string{}
		for _, er := range errors {
			elist = append(elist, er.Error())
		}
		estrs, _ := json.Marshal(elist)
		texcute.Response = string(estrs)
	} else {
		texcute.Response = util.Unicode2utf8(body)
	}

	texcute.Code = code
	texcute.EndTime = util.FormatTimeMill()
	texcute.Duration = time.Now().UnixMilli() - startDuration

	_ = ts.AddExcute(ctx, texcute)

	if code != 200 {
		u, err := ts.UserService.GetById(ctx, t.UserId)
		if err != nil {
			return
		}
		if u.Mail != "" {
			go ts.TasktAlarm(ctx, u.Mail, "错误报警")
		}
	}
}

func (t *taskService) texcutePageByOption(ctx context.Context, opt *types.TaskExcuteOption) (*types.Page[*model.TaskExcute], error) {
	var total int64
	var tasks []*model.TaskExcute

	sess := t.DB.WithContext(ctx)

	offset := (opt.Current - 1) * opt.PageSize
	if offset < 0 {
		offset = 0
	}
	q := sess.Model(&model.Task{})

	if opt.Code != 0 {
		if opt.Code == 1 {
			q = q.Where("code = 200")
		} else {
			q = q.Where("code != 200")
		}
	}

	if opt.TaskGroup != "" {
		q = q.Where("taskGroup LIKE ?", "%"+opt.TaskGroup+"%")
	}

	if opt.TaskName != "" {
		q = q.Where("`taskName` LIKE ?", "%"+opt.TaskName+"%")
	}

	if opt.UserId != 0 {
		q = q.Where("user_id = ?", opt.UserId)
	}

	if opt.TaskId != 0 {
		q = q.Where("task_id = ?", opt.TaskId)
	}

	if opt.StartTime != "" {
		q = q.Where("create_time > ?", opt.StartTime)
	}

	if opt.EndTime != "" {
		q = q.Where("create_time < ?", opt.EndTime)
	}

	q.Count(&total)
	q.Limit(opt.PageSize).Offset(offset).Order("id desc").Find(&tasks)

	pg := &types.Page[*model.TaskExcute]{
		Total: total,
		List:  tasks,
	}

	return pg, nil
}
