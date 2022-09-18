package subscribe

import (
	"encoding/json"
	"fmt"
	"github.com/eininst/flog"
	"github.com/eininst/ninja"
	"github.com/eininst/rlock"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/util"
	"github.com/gofiber/fiber/v2"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"strconv"
	"sync"
	"time"
)

func init() {
	ninja.Provide(&TaskSubscribe{
		TaskMap: make(map[int64]cron.EntryID),
	})
}

const (
	GET    = "get"
	POST   = "post"
	PUT    = "put"
	DELETE = "delete"
)

type TaskSubscribe struct {
	TaskMap map[int64]cron.EntryID
	mux     *sync.Mutex
	RsCli   rs.Client    `inject:""`
	Rlock   *rlock.Rlock `inject:""`
	DB      *gorm.DB     `inject:""`
	CronCli *cron.Cron   `inject:""`

	TaskService task.TaskService `inject:""`
}

func (ts *TaskSubscribe) Register(ctx *rs.Context) {
	defer ctx.Ack()

	task_id := ctx.Msg.Values["task_id"]
	tid, err := strconv.ParseInt(task_id.(string), 10, 64)
	if err != nil {
		return
	}

	var t model.SchedulerTask
	ts.DB.WithContext(ctx).First(&t, tid)
	if t.Id == 0 {
		return
	}

	lockName := fmt.Sprintf("task_run:%v", t.Id)

	eid, er := ts.CronCli.AddFunc(t.Spec, func() {
		ok, cancel := ts.Rlock.TryAcquire(lockName, time.Second)
		defer cancel()
		if !ok {
			return
		}

		tstr, er := json.Marshal(t)
		if er != nil {
			flog.Error(er)
			return
		}
		er = ts.RsCli.Send("task_request", rs.H{
			"info": string(tstr),
		})
		if er != nil {
			flog.Error(er)
		}
	})
	if er != nil {
		flog.Error(er)
		return
	}
	ts.TaskMap[t.Id] = eid
}

func (ts *TaskSubscribe) Stop(ctx *rs.Context) {
	defer ctx.Ack()

	task_id := ctx.Msg.Values["task_id"]
	tid, err := strconv.ParseInt(task_id.(string), 10, 64)
	if err != nil {
		return
	}

	if eid, ok := ts.TaskMap[tid]; ok {
		ts.CronCli.Remove(eid)
		delete(ts.TaskMap, tid)
	}
}

func (ts *TaskSubscribe) Request(ctx *rs.Context) {
	defer ctx.Ack()

	var t model.SchedulerTask
	info := ctx.Msg.Values["info"].(string)
	er := json.Unmarshal([]byte(info), &t)
	if er != nil {
		flog.Error(er)
		return
	}

	texcute := &model.SchedulerTaskExcute{
		UserId:    t.UserId,
		TaskId:    t.Id,
		TaskName:  t.Name,
		TaskUrl:   t.Url,
		TaskObj:   info,
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

	code, body, errors := agt.Bytes()
	if len(errors) > 0 {
		elist := []string{}
		for _, er := range errors {
			elist = append(elist, er.Error())
		}
		estrs, _ := json.Marshal(elist)
		texcute.Response = string(estrs)
	} else {
		texcute.Response = string(body)
	}

	texcute.Code = code

	texcute.EndTime = util.FormatTimeMill()
	texcute.Duration = time.Now().UnixMilli() - startDuration

	_ = ts.TaskService.AddExcute(ctx, texcute)
}
