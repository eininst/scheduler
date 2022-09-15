package api

import (
	"github.com/eininst/flog"
	"github.com/eininst/go-jwt"
	"github.com/eininst/ninja"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"time"
)

func init() {
	ninja.Provide(new(Sapi))
}

type Sapi struct {
	Jwt         *jwt.Jwt         `inject:""`
	UserService user.UserService `inject:""`
}

func (a *Sapi) Index(c *fiber.Ctx) error {
	if c.Path() != "/login" {
		token := c.Cookies("access_token")
		if token == "" {
			return c.Redirect("/login", http.StatusTemporaryRedirect)
		}

		var u model.SchedulerUser
		er := a.Jwt.ParseToken(token, &u)
		if er != nil {
			return er
		}

		return c.Render("index", fiber.Map{
			"user":   u,
			"assets": configs.Get("assets"),
		})
	}

	return c.Render("index", fiber.Map{
		"assets": configs.Get("assets"),
	})
}

func (a *Sapi) Login(c *fiber.Ctx) error {
	var u struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&u); err != nil {
		flog.Error(err)
		return service.NewServiceError("参数错误")
	}
	r, er := a.UserService.Login(c.Context(), u.Username, u.Password)
	if er != nil {
		return er
	}

	dur := time.Hour * 72
	token := a.Jwt.CreateToken(r, dur)
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(dur),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)

	return c.JSON(r)
}

func (a *Sapi) Logout(c *fiber.Ctx) error {
	cookie := new(fiber.Cookie)
	cookie.Name = "access_token"
	cookie.Value = "deleted"
	cookie.HTTPOnly = true
	cookie.Secure = true
	cookie.Expires = time.Now().Add(-3 * time.Second)
	c.Cookie(cookie)

	return c.SendStatus(http.StatusOK)
}

func (a *Sapi) TaskAdd(c *fiber.Ctx) error {
	return nil
}
