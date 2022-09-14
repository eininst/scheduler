package api

import (
	"github.com/eininst/scheduler/internal/common/jwt"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Jwt  *jwt.Jwt   `inject:""`
	App  *fiber.App `inject:""`
	Sapi *Sapi      `inject:""`
}

func (r *Router) RequireLogin(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return service.NewServiceError("用户未登陆")
	}
	var user model.SchedulerUser
	er := r.Jwt.ParseToken(token, &user)
	if er != nil {
		return er
	}
	c.Locals("user", user)
	return c.Next()
}

func (r *Router) Init() {
	r.App.Post("/api/login", r.Sapi.Login)
	r.App.Post("/api/logout", r.Sapi.Logout)

	//g := r.App.Group("/api/u", r.RequireLogin)
	r.App.Get("/*", r.Sapi.Index)
}
