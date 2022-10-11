package task

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/eininst/flog"
	"github.com/eininst/rlock"
	"github.com/eininst/rs"
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

	GET    = "GET"
	POST   = "POST"
	PUT    = "PUT"
	DELETE = "DELETE"
)

var SORT_MAP = map[string]string{
	"createTime": "create_time",
	"start_time": "start_time",
	"duration":   "duration",
}

type TaskService interface {
	RunTask(ctx context.Context)

	Add(ctx context.Context, task *model.Task) error
	Update(ctx context.Context, userId int64, task *model.Task) error
	PageByOption(ctx context.Context, opt *types.TaskOption) (*types.Page[*types.TaskDTO], error)
	Start(ctx context.Context, userId int64, id int64) error
	Stop(ctx context.Context, userId int64, id int64) error
	Del(ctx context.Context, userId int64, id int64) error

	DoOnce(ctx context.Context, uid int64, id int64) error
	RunTaskById(ctx context.Context, taskId int64) error
	DelEntry(ctx context.Context, taskId int64)

	UpdateUser(ctx context.Context, tcu *types.TaskChangeUser) (int64, error)
	StartBatch(ctx context.Context, userId int64, tbatch *types.TaskBatch) (int64, error)
	StopBatch(ctx context.Context, userId int64, tbatch *types.TaskBatch) (int64, error)
	DelBatch(ctx context.Context, userId int64, tbatch *types.TaskBatch) (int64, error)

	AddExcute(ctx context.Context, taskExcute *model.TaskExcute) error

	ExcutePageByOption(ctx context.Context, opt *types.TaskExcuteOption) (*types.Page[*types.TaskExcuteDTO], error)

	CleanLog(ctx context.Context, day int64)

	Dashboard(ctx context.Context) *types.Dashboard
}
type taskService struct {
	parse       cron.Parser
	mux         *sync.Mutex
	TaskMap     map[int64]cron.EntryID
	DB          *gorm.DB         `inject:""`
	Rcli        *redis.Client    `inject:""`
	CronCli     *cron.Cron       `inject:""`
	Rlock       *rlock.Rlock     `inject:""`
	RsClient    rs.Client        `inject:""`
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
	_, er := t.parse.Parse(task.Cron)
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

func (t *taskService) Update(ctx context.Context, userId int64, task *model.Task) error {
	session := t.DB.WithContext(ctx)

	var _task model.Task
	session.First(&_task, task.Id)
	if _task.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	var u model.User
	session.First(&u, userId)
	if u.Role != user.ROLE_ADMIN {
		if u.Id != _task.UserId {
			return service.NewServiceError("修改失败, 权限不足")
		}
	}

	if _task.Status == STATUS_RUN {
		return service.NewServiceError("任务处于运行中，修改失败")
	}

	_, er := t.parse.Parse(task.Cron)
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
		Cron:        task.Cron,
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

func (t *taskService) DoOnce(ctx context.Context, uid int64, id int64) error {
	session := t.DB.WithContext(ctx)

	var u model.User
	session.First(&u, uid)
	if u.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	var task model.Task
	session.First(&task, id)
	if task.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	if u.Role != user.ROLE_ADMIN {
		if task.Id != u.Id {
			return service.NewServiceError("权限不足，执行失败")
		}
	}

	t.do(ctx, &task)

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

func (t *taskService) StartBatch(ctx context.Context, userId int64, tbatch *types.TaskBatch) (int64, error) {
	if len(tbatch.TaskIds) == 0 {
		return 0, service.NewServiceError("请选择任务")
	}

	count := int64(0)
	for _, id := range tbatch.TaskIds {
		er := t.Start(ctx, userId, id)
		if er != nil {
			continue
		}
		count += 1
	}
	return count, nil
}

func (t *taskService) StopBatch(ctx context.Context, userId int64, tbatch *types.TaskBatch) (int64, error) {
	if len(tbatch.TaskIds) == 0 {
		return 0, service.NewServiceError("请选择任务")
	}

	count := int64(0)
	for _, id := range tbatch.TaskIds {
		er := t.Stop(ctx, userId, id)
		if er != nil {
			continue
		}
		count += 1
	}
	return count, nil
}

func (t *taskService) DelBatch(ctx context.Context, userId int64, tbatch *types.TaskBatch) (int64, error) {
	if len(tbatch.TaskIds) == 0 {
		return 0, service.NewServiceError("请选择任务")
	}

	count := int64(0)
	for _, id := range tbatch.TaskIds {
		er := t.Del(ctx, userId, id)
		if er != nil {
			continue
		}
		count += 1
	}
	return count, nil
}

func (ts *taskService) Start(ctx context.Context, userId int64, id int64) error {
	ok, cancel := ts.Rlock.TryAcquire(fmt.Sprintf("task:%v", id), time.Second*10)
	defer cancel()
	if !ok {
		return service.NewServiceError("启动失败, 当前操作有其他人正在进行操作")
	}

	sess := ts.DB.WithContext(ctx)

	var t model.Task
	sess.First(&t, id)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	var u model.User
	sess.First(&u, userId)

	if u.Role != user.ROLE_ADMIN {
		if u.Id != t.UserId {
			return service.NewServiceError("启动失败, 权限不足")
		}
	}

	if t.Status == STATUS_RUN {
		return service.NewServiceError("任务已启动，请勿重复操作")
	}
	er := sess.Model(&t).Update("status", STATUS_RUN).Error
	if er != nil {
		return service.NewServiceError("启动失败")
	}

	er = ts.RsClient.Send("task_run", rs.H{
		"task_id": t.Id,
	})
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) Stop(ctx context.Context, userId int64, id int64) error {
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

	var u model.User
	sess.First(&u, userId)

	if u.Role != user.ROLE_ADMIN {
		if u.Id != t.UserId {
			return service.NewServiceError("停止失败, 权限不足")
		}
	}

	if t.Status == STATUS_STOP {
		return service.NewServiceError("任务已停止，请勿重复操作")
	}

	er := sess.Model(&t).Update("status", STATUS_STOP).Error
	if er != nil {
		return service.NewServiceError("停止失败")
	}

	er = ts.RsClient.Send("task_stop", rs.H{
		"task_id": t.Id,
	})
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) Del(ctx context.Context, userId int64, id int64) error {
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

	var u model.User
	sess.First(&u, userId)

	if u.Role != user.ROLE_ADMIN {
		if u.Id != t.UserId {
			return service.NewServiceError("删除失败, 权限不足")
		}
	}

	er := sess.Delete(&t).Error
	if er != nil {
		return service.NewServiceError("删除失败")
	}

	er = ts.RsClient.Send("task_stop", rs.H{
		"task_id": t.Id,
	})
	if er != nil {
		return er
	}

	return nil
}

func (ts *taskService) runTask(ctx context.Context, t *model.Task) error {
	eid, er := ts.CronCli.AddFunc(t.Cron, func() {
		ts.do(ctx, t)
	})
	if er != nil {
		return er
	}

	ts.addEntry(ctx, t.Id, eid)
	return nil
}

func (ts *taskService) RunTaskById(ctx context.Context, taskId int64) error {
	var t model.Task
	ts.DB.WithContext(ctx).First(&t, taskId)
	if t.Id == 0 {
		return service.ERROR_DATA_NOT_FOUND
	}

	eid, er := ts.CronCli.AddFunc(t.Cron, func() {
		ts.do(ctx, &t)
	})
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

	q = q.Count(&total)

	q = q.Limit(opt.PageSize).Offset(offset)

	if opt.Sort != "" {
		sort, ok := SORT_MAP[opt.Sort]
		if ok {
			if opt.Dir == "ascend" {
				q = q.Order(fmt.Sprintf("%s", sort))
			} else {
				q = q.Order(fmt.Sprintf("%s desc", sort))
			}
		}
	} else {
		q = q.Order("id desc")
	}

	q.Find(&tasks)
	//sort
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

func (ts *taskService) DelEntry(ctx context.Context, taskId int64) {
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
		if er != nil {
			flog.Fatal(er)
		}
	}
}

func (ts *taskService) CleanLog(ctx context.Context, day int64) {
	f := time.Now().Add(-time.Hour * 24 * time.Duration(day)).Format("2006-01-02 15:04:05")
	te := &model.TaskExcute{}
	dsql := fmt.Sprintf("delete from %s where create_time < ?", te.TableName())
	ts.DB.WithContext(ctx).Exec(dsql, f)
}

func (ts *taskService) do(ctx context.Context, t *model.Task) {
	lockName := fmt.Sprintf("task_do:%v", t.Id)
	ok, cancel := ts.Rlock.TryAcquire(lockName, time.Second*time.Duration(t.Timeout))
	defer cancel()
	if !ok {
		return
	}
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
		sbody := body
		if len(sbody) > 500 {
			sbody = sbody[0:500]
		}
		texcute.Response = sbody
	}

	texcute.Code = code
	texcute.EndTime = util.FormatTimeMill()
	texcute.Duration = time.Now().UnixMilli() - startDuration

	texcuteStr, _ := json.Marshal(texcute)

	_ = ts.RsClient.Send("cron_task_log", rs.H{
		"data": string(texcuteStr),
	})
}

func (t *taskService) ExcutePageByOption(ctx context.Context, opt *types.TaskExcuteOption) (*types.Page[*types.TaskExcuteDTO], error) {
	var total int64
	var tasks []*model.TaskExcute

	sess := t.DB.WithContext(ctx)

	offset := (opt.Current - 1) * opt.PageSize
	if offset < 0 {
		offset = 0
	}
	q := sess.Model(&model.TaskExcute{})

	if opt.Status != "" {
		if opt.Status == "ok" {
			q = q.Where("code = 200")
		} else {
			q = q.Where("code != 200")
		}
	}

	if opt.TaskGroup != "" {
		q = q.Where("task_group LIKE ?", "%"+opt.TaskGroup+"%")
	}

	if opt.TaskName != "" {
		q = q.Where("`task_name` LIKE ?", "%"+opt.TaskName+"%")
	}

	if opt.UserId != 0 {
		q = q.Where("user_id = ?", opt.UserId)
	}

	if opt.TaskId != 0 {
		q = q.Where("task_id = ?", opt.TaskId)
	}

	if opt.StartTime != "" {
		q = q.Where("start_time > ?", fmt.Sprintf("%s 00:00:0", opt.StartTime))
	}

	if opt.EndTime != "" {
		q = q.Where("start_time < ?", fmt.Sprintf("%s 23:59:59", opt.EndTime))
	}
	if opt.Duration != 0 {
		q = q.Where("duration > ? ", opt.Duration)
	}

	q.Count(&total)

	q = q.Limit(opt.PageSize).Offset(offset)
	//q.Limit(opt.PageSize).Offset(offset).Order("id desc").Find(&tasks)

	if opt.Sort != "" {
		sort, ok := SORT_MAP[opt.Sort]
		if ok {
			if opt.Dir == "ascend" {
				q = q.Order(fmt.Sprintf("%s", sort))
			} else {
				q = q.Order(fmt.Sprintf("%s desc", sort))
			}
		}
	} else {
		q = q.Order("id desc")
	}
	q.Find(&tasks)

	taskDtos := []*types.TaskExcuteDTO{}

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
		}
	}

	pg := &types.Page[*types.TaskExcuteDTO]{
		Total: total,
		List:  taskDtos,
	}

	return pg, nil
}

func (t *taskService) Dashboard(ctx context.Context) *types.Dashboard {
	sess := t.DB.WithContext(ctx)

	var taskCount int64
	var taskRunCount int64
	var schedulerCount int64

	sess.Model(&model.Task{}).Count(&taskCount)
	sess.Model(&model.Task{}).Where("status = ?", STATUS_RUN).Count(&taskRunCount)

	sess.Model(&model.TaskExcute{}).Count(&schedulerCount)

	sql := `select 
	left(create_time, 10) 'date',
	code,
	count(*) 'count'
	from scheduler_task_excute where create_time>? group by left(create_time, 10), code;`

	var chart []*types.DashboardChart
	tm := time.Now().Add(-time.Hour * 24 * time.Duration(10)).Format("2006-01-02")

	sess.Raw(sql, fmt.Sprintf("%s 00:00:00", tm)).Find(&chart)

	return &types.Dashboard{
		TaskCount:      taskCount,
		TaskRunCount:   taskRunCount,
		SchedulerCount: schedulerCount,
		Chart:          chart,
		StartTime:      tm,
	}
}
