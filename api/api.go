package api

import (
	"github.com/eininst/flog"
	"github.com/eininst/go-jwt"
	"github.com/eininst/ninja"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/eininst/scheduler/internal/types"
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
	TaskService task.TaskService `inject:""`
}

func (a *Sapi) Index(c *fiber.Ctx) error {
	if c.Path() != "/login" {
		token := c.Cookies("access_token")
		if token == "" {
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

	return nil
}

func (a *Sapi) UserAdd(c *fiber.Ctx) error {
	var u model.User
	er := c.BodyParser(&u)
	if er != nil {
		return er
	}

	return a.UserService.Add(c.Context(), &u)
}

func (a *Sapi) UserList(c *fiber.Ctx) error {
	users, er := a.UserService.List(c.Context())
	if er != nil {
		return er
	}
	return c.JSON(users)
}

func (a *Sapi) TaskAdd(c *fiber.Ctx) error {
	var t model.Task
	er := c.BodyParser(&t)
	if er != nil {
		return er
	}
	uid := c.Locals("userId").(int64)
	t.UserId = uid

	return a.TaskService.Add(c.Context(), &t)
}

func (a *Sapi) TaskUpdate(c *fiber.Ctx) error {
	var t model.Task
	er := c.BodyParser(&t)
	if er != nil {
		return er
	}
	//uid := c.Locals("userId").(int64)
	//t.UserId = uid

	return a.TaskService.Update(c.Context(), &t)
}

func (a *Sapi) TaskPage(c *fiber.Ctx) error {
	var opt types.TaskOption

	er := c.QueryParser(&opt)
	if er != nil {
		return er
	}

	r, er := a.TaskService.PageByOption(c.Context(), &opt)
	if er != nil {
		return er
	}

	return c.JSON(r)
}

func (a *Sapi) TaskStart(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	return a.TaskService.Start(c.Context(), int64(id))
}

func (a *Sapi) TaskStop(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	return a.TaskService.Stop(c.Context(), int64(id))
}

func (a *Sapi) TaskDel(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	return a.TaskService.Del(c.Context(), int64(id))
}

func (a *Sapi) TaskUpdateUser(c *fiber.Ctx) error {
	var changeUser types.TaskChangeUser

	er := c.BodyParser(&changeUser)
	if er != nil {
		return er
	}
	count, er := a.TaskService.UpdateUser(c.Context(), &changeUser)
	if er != nil {
		return er
	}
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (a *Sapi) StartBatch(c *fiber.Ctx) error {
	var tbatch types.TaskBatch

	er := c.BodyParser(&tbatch)
	if er != nil {
		return er
	}
	count, er := a.TaskService.StartBatch(c.Context(), &tbatch)
	if er != nil {
		return er
	}
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (a *Sapi) StopBatch(c *fiber.Ctx) error {
	var tbatch types.TaskBatch

	er := c.BodyParser(&tbatch)
	if er != nil {
		return er
	}
	count, er := a.TaskService.StopBatch(c.Context(), &tbatch)
	if er != nil {
		return er
	}
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (a *Sapi) DelBatch(c *fiber.Ctx) error {
	var tbatch types.TaskBatch

	er := c.BodyParser(&tbatch)
	if er != nil {
		return er
	}
	count, er := a.TaskService.DelBatch(c.Context(), &tbatch)
	if er != nil {
		return er
	}
	return c.JSON(fiber.Map{
		"count": count,
	})
}
