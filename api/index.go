package api

import (
	"encoding/json"
	"fmt"
	"github.com/eininst/go-jwt"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/inject"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"html/template"
	"net/http"
	"time"
)

func init() {
	inject.Provide(new(IndexApi))
}

type IndexApi struct {
	Jwt         *jwt.Jwt         `inject:""`
	UserService user.UserService `inject:""`
}

func (a *IndexApi) Index() fiber.Handler {
	index, err := template.New("index.html").Parse(tpl)
	if err != nil {
		panic(fmt.Errorf("fiber: swagger middleware error -> %w", err))
	}

	defaultWebConfig := &WebConfig{
		Assets: "assets",
		Title:  "Scheduler",
		Desc:   "简单，开箱即用的定时任务平台",
		Logo:   "",
		Avatar: "",
	}
	var wcfg WebConfig

	webConfigStr := configs.Get("web").String()
	if webConfigStr == "" {
		wcfg = *defaultWebConfig
	} else {
		er := json.Unmarshal([]byte(webConfigStr), &wcfg)
		if er != nil {
			wcfg = *defaultWebConfig
		}
	}
	if wcfg.Assets == "" {
		wcfg.Assets = defaultWebConfig.Assets
	}
	if wcfg.Title == "" {
		wcfg.Title = defaultWebConfig.Title
	}
	if wcfg.Avatar == "" {
		wcfg.Avatar = defaultWebConfig.Avatar
	}

	return func(c *fiber.Ctx) error {
		c.Type("html")

		if c.Path() == "/init" {
			count, _ := a.UserService.Count(c.Context())
			if count == 0 {
				return index.Execute(c, fiber.Map{
					"config": wcfg,
				})
			}
			token := c.Cookies("access_token")
			if token == "" {
				return c.Redirect("/login", http.StatusTemporaryRedirect)
			}
			return c.Redirect("/", http.StatusTemporaryRedirect)
		}
		if c.Path() != "/login" && c.Path() != "/init" {
			token := c.Cookies("access_token")
			if token == "" {
				count, _ := a.UserService.Count(c.Context())
				if count == 0 {
					return c.Redirect("/init", http.StatusTemporaryRedirect)
				}
				return c.Redirect("/login", http.StatusTemporaryRedirect)
			}

			var u model.User
			er := a.Jwt.ParseToken(token, &u)
			if er != nil {
				return er
			}

			nu, er := a.UserService.GetById(c.Context(), u.Id)
			if er == nil {
				u = *nu
			}
			if nu.Status != user.STATUS_OK {
				cookie := new(fiber.Cookie)
				cookie.Name = "access_token"
				cookie.Value = "deleted"
				cookie.HTTPOnly = true
				cookie.Secure = true
				cookie.Expires = time.Now().Add(-3 * time.Second)
				c.Cookie(cookie)
				return c.Redirect("/login", http.StatusTemporaryRedirect)
			}

			return index.Execute(c, fiber.Map{
				"user":   u,
				"config": wcfg,
			})
		}

		return index.Execute(c, fiber.Map{
			"config": wcfg,
		})
	}
}

const tpl = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8"/>
  <meta
    name="viewport"
    content="width=device-width, initial-scale=1, maximum-scale=1, minimum-scale=1, user-scalable=no"
  />
  <link rel="stylesheet" href="{{.config.Assets}}/umi.css"/>
  <script>
    window.routerBase = "/";
    window.config = {
      title: "{{.config.Title}}",
      desc: "{{.config.Desc}}",
      avatar: "{{.config.Avatar}}",
      logo: "{{.config.Logo}}"
    }

    window.role = "{{.user.Role}}"
    window.userName = "{{.user.Name}}"
  </script>
</head>
<body>
<div id="root"></div>

<script src="{{.config.Assets}}/umi.js"></script>
</body>
</html>
`
