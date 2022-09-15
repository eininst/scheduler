package main

import (
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/flog"
	"github.com/eininst/ninja"
	"github.com/eininst/scheduler/api"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
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

	ninja.Install(new(api.Router), app)

	grace.Listen(app, ":8999")
}
