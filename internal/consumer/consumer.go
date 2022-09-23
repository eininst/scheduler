package consumer

import (
	"encoding/json"
	"github.com/eininst/flog"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
)

type Consumer struct {
	RsClient    rs.Client        `inject:""`
	TaskService task.TaskService `inject:""`
	UserService user.UserService `inject:""`
}

func New() *Consumer {
	return &Consumer{}
}

func (c *Consumer) Init() {
	c.RsClient.Receive(rs.Rctx{
		Stream: "cron_task_log",
		Handler: func(ctx *rs.Context) {
			defer ctx.Ack()
			var excute model.TaskExcute

			dataStr := ctx.Msg.Values["data"].(string)

			er := json.Unmarshal([]byte(dataStr), &excute)
			if er != nil {
				flog.Error(er)
			}
			er = c.TaskService.AddExcute(ctx, &excute)
			if er != nil {
				flog.Error(er)
			}
			if excute.Code != 200 {
				_ = c.RsClient.Send("cron_task_alarm", rs.H{
					"uid": "xx",
					"msg": "msg",
				})
			}
		},
	})

	c.RsClient.Receive(rs.Rctx{
		Stream: "cron_task_alarm",
		Handler: func(ctx *rs.Context) {
			defer ctx.Ack()
			flog.Info(ctx.Msg.Values)
		},
	})
}
