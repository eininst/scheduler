package main

import (
	"context"
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/flog"
	"github.com/eininst/ninja"
	"github.com/eininst/scheduler/api"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/consumer"
	"github.com/eininst/scheduler/internal/conf"
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

}

func setup() {
	var s struct {
		TaskService task.TaskService `inject:""`
	}
	ninja.Populate(&s)
	s.TaskService.Run(context.Background())
}

func main() {
	//cron
	cronCli := cron.New(cron.WithSeconds())
	ninja.Populate(cronCli)
	cronCli.Start()

	// consumer
	var csm consumer.Consumer
	ninja.Install(&csm)
	go csm.Cli.Listen()

	//task run...
	setup()

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

	//grace stop
	cronCli.Stop()
	csm.Cli.Shutdown()
}
