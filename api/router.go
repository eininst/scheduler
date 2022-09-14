package api

import (
	"github.com/gofiber/fiber/v2"
)

type Router struct {
	App  *fiber.App `inject:""`
	Sapi *Sapi      `inject:""`
}

func (r *Router) Init() {
	r.App.Post("/api/login", r.Sapi.Login)
	r.App.Post("/api/logout", r.Sapi.Logout)
	r.App.Get("/*", r.Sapi.Index)
}
