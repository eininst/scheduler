package scheduler

import (
	"context"
	"fmt"
	"github.com/eininst/flog"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/api"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/consumer"
	"github.com/eininst/scheduler/internal/inject"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html"
	"github.com/robfig/cron/v3"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type App interface {
	Listen(cfg ...Config)
}

type app struct {
	TaskService task.TaskService   `inject:""`
	CronCli     *cron.Cron         `inject:""`
	Consumer    *consumer.Consumer `inject:""`
	RsClient    rs.Client          `inject:""`
}

type Config struct {
	Port string
	Sig  os.Signal
}

func New(cfgPath string) App {
	configs.SetConfig(cfgPath)
	conf.Inject()
	inject.Provide(consumer.New())
	app := &app{}
	inject.Provide(app)
	inject.Populate()
	return app
}

func (a *app) cronStart() *cron.Cron {
	a.CronCli.Start()
	a.TaskService.RunTask(context.Background())

	go func() {
		for {
			time.Sleep(time.Minute * 10)
			a.TaskService.CleanLog(context.TODO(), 10)
		}
	}()
	return a.CronCli
}

func (a *app) Listen(config ...Config) {
	port := configs.Get("port").String()
	var sig os.Signal
	if len(config) > 0 {
		cfg := config[0]
		port = cfg.Port
		if cfg.Sig != nil {
			sig = cfg.Sig
		}
	}

	if sig == nil {
		sig = syscall.SIGTERM
	}

	flog.Info(sig)

	if port == "" {
		flog.Fatal("port is required in config.yaml")
	}
	if !strings.HasPrefix(port, ":") {
		port = fmt.Sprintf(":%s", port)
	}

	a.CronCli.Start()
	a.TaskService.RunTask(context.Background())
	retainDay := configs.Get("log", "retain").Int()
	interval := configs.Get("log", "interval").Int()
	if retainDay == 0 {
		retainDay = 10
	}
	if interval == 0 {
		interval = 60 * 5
	}
	go func() {
		for {
			time.Sleep(time.Second * time.Duration(interval))
			a.TaskService.CleanLog(context.TODO(), retainDay)
		}
	}()

	a.Consumer.Init()
	go a.RsClient.Listen()

	//app
	engine := html.New("./website/views", ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		Prefork:      false,
		ReadTimeout:  time.Second * 10,
		ErrorHandler: service.ErrorHandler,
	})
	app.Static("/assets", "./website/dist")

	title := configs.Get("web", "title").String()
	if title == "" {
		title = "Scheduler"
	}
	app.Get("/metrics", monitor.New(monitor.Config{
		Title: title,
	}))

	inject.Install(new(api.Router), app)

	go func() { _ = app.Listen(port) }()

	quit := make(chan os.Signal)
	signal.Notify(quit, sig)
	<-quit
	a.CronCli.Stop()
	a.Consumer.RsClient.Shutdown()
	_ = app.Shutdown()

	log.Println("Shutdown Server ...")
}
