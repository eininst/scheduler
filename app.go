package scheduler

import (
	"context"
	"github.com/eininst/ninja"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/api"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/consumer"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/template/html"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App interface {
	Listen(port string)
}

type app struct {
	TaskService task.TaskService   `inject:""`
	CronCli     *cron.Cron         `inject:""`
	Consumer    *consumer.Consumer `inject:""`
	RsClient    rs.Client          `inject:""`
}

//var Assets = "/assets"

//var defaultSecretKey = "b956160659554dbcb0ae65e2f7d5de14"

type Config struct {
	RedisClient *redis.Client
	DB          *gorm.DB
	SecretKey   string
	LogWorker   int64
}

func New(cfgPath string) App {
	configs.SetConfig(cfgPath)
	conf.Inject()
	ninja.Provide(consumer.New())
	app := &app{}
	ninja.Provide(app)
	ninja.Populate()
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
func (a *app) Listen(port string) {
	a.CronCli.Start()
	a.TaskService.RunTask(context.Background())

	go func() {
		for {
			time.Sleep(time.Minute * 10)
			a.TaskService.CleanLog(context.TODO(), 10)
		}
	}()

	a.Consumer.Init()
	go a.RsClient.Listen()

	//app
	engine := html.New("./web/views", ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		Prefork:      false,
		ReadTimeout:  time.Second * 10,
		ErrorHandler: service.ErrorHandler,
	})
	app.Static("/assets", "./web/dist")
	app.Get("/metrics", monitor.New(monitor.Config{
		Title: "Scheduler",
	}))

	ninja.Install(new(api.Router), app)

	go func() { _ = app.Listen(port) }()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM)
	<-quit
	a.CronCli.Stop()
	a.Consumer.RsClient.Shutdown()
	_ = app.Shutdown()

	log.Println("Shutdown Server ...")
}
