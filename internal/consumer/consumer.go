package consumer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/eininst/flog"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service/mail"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/valyala/fasttemplate"
	"strconv"
	"strings"
)

type Consumer struct {
	RsClient    rs.Client        `inject:""`
	TaskService task.TaskService `inject:""`
	UserService user.UserService `inject:""`
	MailService mail.MailService `inject:""`
}

func New() *Consumer {
	return &Consumer{}
}

func (c *Consumer) runTask(ctx *rs.Context) {
	defer ctx.Ack()
	taskIdStr := ctx.Msg.Values["task_id"]

	taskId, err := strconv.ParseInt(taskIdStr.(string), 10, 64)
	if err != nil {
		return
	}

	er := c.TaskService.RunTaskById(ctx, taskId)
	if er != nil {
		flog.Error(er)
	}
}

func (c *Consumer) stopTask(ctx *rs.Context) {
	defer ctx.Ack()
	taskIdStr := ctx.Msg.Values["task_id"]
	taskId, err := strconv.ParseInt(taskIdStr.(string), 10, 64)
	if err != nil {
		return
	}

	c.TaskService.DelEntry(ctx, taskId)
}

func (c *Consumer) Init() {
	taskLogWork := configs.Get("log", "  work").Int()
	if taskLogWork == 0 {
		taskLogWork = 5
	}
	c.RsClient.Receive(rs.Rctx{
		Stream:  "task_run",
		Work:    rs.Int(1024),
		Group:   utils.UUIDv4(),
		Handler: c.runTask,
	})

	c.RsClient.Receive(rs.Rctx{
		Stream:  "task_stop",
		Work:    rs.Int(1024),
		Group:   utils.UUIDv4(),
		Handler: c.stopTask,
	})

	c.RsClient.Receive(rs.Rctx{
		Stream: "cron_task_log",
		Work:   rs.Int(int(taskLogWork)),
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
			if excute.Code != 200 && c.MailService.IsConfig() {
				_ = c.RsClient.Send("cron_task_alarm", rs.H{
					"data": dataStr,
				})
			}
		},
	})

	mailWork := configs.Get("mail", "  work").Int()
	if mailWork == 0 {
		mailWork = 5
	}

	c.RsClient.Receive(rs.Rctx{
		Stream: "cron_task_alarm",
		Work:   rs.Int(int(mailWork)),
		Handler: func(ctx *rs.Context) {
			defer ctx.Ack()

			var excute model.TaskExcute
			dataStr := ctx.Msg.Values["data"].(string)

			er := json.Unmarshal([]byte(dataStr), &excute)
			if er != nil {
				flog.Error(er)
				return
			}

			u, er := c.UserService.GetById(ctx, excute.UserId)
			if er != nil {
				flog.Error(er)
				return
			}
			if u.Mail == "" {
				return
			}

			var tk model.Task
			er = json.Unmarshal([]byte(excute.TaskObj), &tk)
			if er != nil {
				flog.Error(er)
				return
			}

			subject := fmt.Sprintf(`任务：%s, 出错了`, excute.TaskName)

			if json.Valid([]byte(excute.Response)) {
				var out bytes.Buffer
				err := json.Indent(&out, []byte(excute.Response), "", "\t")
				if err == nil {
					excute.Response = out.String()
				}
			}

			t := fasttemplate.New(tpl, "${", "}")
			var parms = map[string]any{
				"url":    excute.TaskUrl,
				"method": strings.ToUpper(tk.Method),
				"rtime":  fmt.Sprintf("%v", excute.Duration),
				"code":   fmt.Sprintf("%v", excute.Code),
				"resp":   excute.Response,
				"body":   "",
			}

			if tk.Method != task.GET {
				parms["body"] = fmt.Sprintf(`
				<p>Content-Type: %s</p>
				<p>Body: %s</p>
				`, strings.ToUpper(tk.Method), tk.Body)
			}
			msg := t.ExecuteString(parms)

			er = c.MailService.Send(u.Mail, subject, msg)
			if er != nil {
				flog.Error(er)
			}
		},
	})
}

const tpl = `
<div>
<p>URL: ${url}</p>
<p>Method: ${method}</p>
${body}
<p>请求耗时: ${rtime}ms</p>
<p>Response Code: ${code}</p>
<p>Response:</p>
<pre>${resp}</pre>
</div>
`
