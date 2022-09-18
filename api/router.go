package api

import (
	"errors"
	"github.com/eininst/go-jwt"
	"github.com/eininst/rs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Jwt   *jwt.Jwt   `inject:""`
	App   *fiber.App `inject:""`
	Sapi  *Sapi      `inject:""`
	RsCli rs.Client  `inject:""`
}

func (r *Router) RequireLogin(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return service.NewServiceError("用户未登陆")
	}
	var user model.SchedulerUser
	er := r.Jwt.ParseToken(token, &user)
	if errors.Is(er, jwt.Expired) {
		return service.NewServiceError("token is expired")
	}

	c.Locals("user", user)
	c.Locals("userId", user.Id)
	return c.Next()
}

func (r *Router) Init() {
	r.App.Get("/test", func(ctx *fiber.Ctx) error {
		return r.RsCli.Send("task_add", rs.H{
			"task_id": int64(123),
		})
	})
	r.App.Post("/api/login", r.Sapi.Login)
	r.App.Post("/api/logout", r.Sapi.Logout)

	g := r.App.Group("/api/u", r.RequireLogin)
	g.Get("/user", r.Sapi.UserList)

	g.Post("/task/add", r.Sapi.TaskAdd)
	g.Put("/task/update", r.Sapi.TaskUpdate)
	g.Get("/task/page", r.Sapi.TaskPage)

	g.Post("/task/start/:id", r.Sapi.TaskStart)
	g.Post("/task/stop/:id", r.Sapi.TaskStop)
	g.Delete("/task/del/:id", r.Sapi.TaskDel)

	r.App.Get("/*", r.Sapi.Index)
}
