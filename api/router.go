package api

import (
	"errors"
	"github.com/eininst/go-jwt"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	Jwt         *jwt.Jwt         `inject:""`
	App         *fiber.App       `inject:""`
	Sapi        *Sapi            `inject:""`
	IndexApi    *IndexApi        `inject:""`
	UserService user.UserService `inject:""`
}

func (r *Router) RequireLogin(c *fiber.Ctx) error {
	token := c.Cookies("access_token")
	if token == "" {
		return service.NewServiceError("用户未登陆")
	}
	var u model.User
	er := r.Jwt.ParseToken(token, &u)
	if errors.Is(er, jwt.Expired) {
		return service.NewServiceError("token is expired")
	}

	c.Locals("user", u)
	c.Locals("userId", u.Id)
	return c.Next()
}

func (r *Router) RequireAdminRole(c *fiber.Ctx) error {
	userId := c.Locals("userId").(int64)
	u, er := r.UserService.GetById(c.Context(), userId)
	if er != nil {
		return er
	}
	if u.Role != user.ROLE_ADMIN {
		return service.NewServiceError("权限不足")
	}
	return c.Next()
}

func (r *Router) Init() {
	r.App.Get("/err", func(ctx *fiber.Ctx) error {
		return service.NewServiceError("error test")
	})
	r.App.Post("/api/init", r.Sapi.Init)
	r.App.Post("/api/login", r.Sapi.Login)
	r.App.Post("/api/logout", r.Sapi.Logout)

	g := r.App.Group("/api/u", r.RequireLogin)

	g.Get("/dashboard", r.Sapi.Dashboard)
	g.Get("/user", r.Sapi.UserList)

	g.Post("/user/add", r.RequireAdminRole, r.Sapi.UserAdd)
	g.Put("/user/update", r.RequireAdminRole, r.Sapi.UserUpdate)
	g.Put("/user/password/reset", r.RequireAdminRole, r.Sapi.UserResetPassword)
	g.Post("/user/enable/:id", r.RequireAdminRole, r.Sapi.UserEnable)
	g.Post("/user/disable/:id", r.RequireAdminRole, r.Sapi.UserDisable)
	g.Delete("/user/del/:id", r.RequireAdminRole, r.Sapi.UserDel)

	g.Post("/task/add", r.Sapi.TaskAdd)
	g.Put("/task/update", r.Sapi.TaskUpdate)
	g.Post("/task/do/:id", r.Sapi.TaskDo)

	g.Get("/task/page", r.Sapi.TaskPage)

	g.Post("/task/start/:id", r.Sapi.TaskStart)
	g.Post("/task/stop/:id", r.Sapi.TaskStop)
	g.Delete("/task/del/:id", r.Sapi.TaskDel)

	g.Post("/task/batch/change/user", r.RequireAdminRole, r.Sapi.TaskUpdateUser)
	g.Post("/task/batch/start", r.RequireAdminRole, r.Sapi.StartBatch)
	g.Post("/task/batch/stop", r.RequireAdminRole, r.Sapi.StopBatch)
	g.Post("/task/batch/del", r.RequireAdminRole, r.Sapi.DelBatch)

	g.Get("/task/excute/page", r.Sapi.ExcutePage)

	r.App.Get("/*", r.IndexApi.Index())
}
