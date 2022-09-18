package main

import (
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/flog"
	"github.com/eininst/ninja"
	"github.com/eininst/scheduler/api"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/consumer"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/robfig/cron/v3"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	logf := "%s[${pid}]%s ${time} ${level} ${path} ${msg}"
	flog.SetFormat(fmt.Sprintf(logf, flog.Cyan, flog.Reset))

	configs.SetConfig("./configs/config.yml")

	conf.Inject()

}

func main() {
	//fiber.MIMEApplicationForm
	engine := html.New("./web/views", ".html")
	app := fiber.New(fiber.Config{
		Views:        engine,
		ReadTimeout:  time.Second * 10,
		ErrorHandler: service.ErrorHandler,
	})
	app.Static("/assets", "./web/dist")

	cronCli := cron.New(cron.WithSeconds())
	ninja.Provide(cronCli)

	ninja.Install(new(api.Router), app)

	var csm consumer.Consumer
	ninja.Install(&csm)

	cronCli.Start()
	go csm.Cli.Listen()

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		cronCli.Stop()
		csm.Cli.Shutdown()
		flog.Info("Graceful Shutdown")
	}()

	grace.Listen(app, ":8999")
}
