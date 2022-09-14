package main

import (
	"fmt"
	grace "github.com/eininst/fiber-prefork-grace"
	"github.com/eininst/flog"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/conf"
	"github.com/eininst/scheduler/internal/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/template/html"
	"net/http"
	"time"
)

func init() {
	logf := "%s[${pid}]%s ${time} ${level} ${path} ${msg}"
	flog.SetFormat(fmt.Sprintf(logf, flog.Cyan, flog.Reset))

	configs.SetConfig("./configs/config.yml")

	conf.Inject()
}

func main() {
	engine := html.New("./web/views", ".html")
	app := fiber.New(fiber.Config{
		Prefork:      false,
		Views:        engine,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
		ErrorHandler: service.ErrorHandler,
	})

	app.Use("assets", filesystem.New(filesystem.Config{
		Root: http.Dir("./web/dist"),
	}))

	app.Get("/*", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"assets": configs.Get("assets"),
		})
	})

	grace.Listen(app, ":8999")
}
