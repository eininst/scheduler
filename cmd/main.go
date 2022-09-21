package main

import (
	"context"
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/flog"
	"github.com/eininst/ninja"
	"github.com/eininst/scheduler/api"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/consumer"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/robfig/cron/v3"
	"time"
)

func init() {
	logf := "%s[${pid}]%s ${time} ${level} ${path} ${msg}"
	flog.SetFormat(fmt.Sprintf(logf, flog.Cyan, flog.Reset))

	configs.SetConfig("./configs/config.yml")

	conf.Inject()
	ninja.Populate()
}

func cronSetup() *cron.Cron {
	var s struct {
		TaskService task.TaskService `inject:""`
		CronCli     *cron.Cron       `inject:""`
	}
	ninja.Populate(&s)

	s.CronCli.Start()
	s.TaskService.RunTask(context.Background())

	go func() {
		for {
			time.Sleep(time.Minute * 10)
			s.TaskService.CleanLog(context.TODO(), 10)
		}
	}()
	return s.CronCli
}

func waitCron(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(time.Millisecond * 10)
		}
	}
}
func main() {
	//task run...
	cronCli := cronSetup()

	var c consumer.Consumer
	ninja.Install(&c)

	go c.RsClient.Listen()

	//app
	engine := html.New("./web/views", ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		ReadTimeout:  time.Second * 10,
		ErrorHandler: service.ErrorHandler,
	})
	app.Static("/assets", "./web/dist")
	ninja.Install(new(api.Router), app)

	//listen
	grace.Listen(app, ":8999")

	//cron grace stop
	ctx := cronCli.Stop()
	waitCron(ctx)

}
