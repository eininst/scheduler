package main

import (
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/flog"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/service"
	"github.com/gofiber/fiber/v2"
	"time"
)

func init() {
	logf := "%s[${pid}]%s ${time} ${level} ${path} ${msg}"
	flog.SetFormat(fmt.Sprintf(logf, flog.Cyan, flog.Reset))

	configs.SetConfig("./configs/helloword.yml")

	conf.Inject()
}

func main() {
	app := fiber.New(fiber.Config{
		Prefork:      false,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: service.ErrorHandler,
	})

	app.Get("/test", func(ctx *fiber.Ctx) error {
		flog.Info("test.1")
		time.Sleep(time.Second * 2)
		return ctx.SendString("ww")
	})

	grace.Listen(app, ":8080")
}
