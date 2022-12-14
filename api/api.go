package api

import (
	"github.com/eininst/flog"
	"github.com/eininst/go-jwt"
	"github.com/eininst/scheduler/internal/inject"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/eininst/scheduler/internal/types"
	"github.com/gofiber/fiber/v2"
	"time"
)

func init() {
	inject.Provide(new(Sapi))
}

type Sapi struct {
	Jwt         *jwt.Jwt         `inject:""`
	UserService user.UserService `inject:""`
	TaskService task.TaskService `inject:""`
}

type WebConfig struct {
	Assets string `json:"assets"`
	Title  string `json:"title"`
	Desc   string `json:"desc"`
	Logo   string `json:"logo"`
	Avatar string `json:"avatar"`
}

func (a *Sapi) Init(c *fiber.Ctx) error {
	count, er := a.UserService.Count(c.Context())
	if er != nil {
		return er
	}
	if count != 0 {
		return service.NewServiceError("初始化设置已完成，请勿重复操作")
	}
	var u model.User
	er = c.BodyParser(&u)
	if er != nil {
		return er
	}

	er = a.UserService.Add(c.Context(), &u)
	if er != nil {
		return er
	}

	dur := time.Hour * 72
	token := a.Jwt.CreateToken(u, dur)
	cookie := fiber.Cookie{
		Name:     "access_token",
		Value:    token,
		Expires:  time.Now().Add(dur),
		HTTPOnly: true,
		Secure:   true,
	}
	c.Cookie(&cookie)
	return c.JSON(u)
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

func (a *Sapi) UserUpdate(c *fiber.Ctx) error {
	var u model.User
	er := c.BodyParser(&u)
	if er != nil {
		return er
	}

	return a.UserService.Update(c.Context(), &u)
}

func (a *Sapi) UserResetPassword(c *fiber.Ctx) error {
	var u model.User
	er := c.BodyParser(&u)
	if er != nil {
		return er
	}

	return a.UserService.ResetPassword(c.Context(), u.Id, u.Password)
}

func (a *Sapi) UserEnable(c *fiber.Ctx) error {
	id, er := c.ParamsInt("id")
	if er != nil {
		return er
	}
	return a.UserService.Enable(c.Context(), int64(id))
}

func (a *Sapi) UserDisable(c *fiber.Ctx) error {
	id, er := c.ParamsInt("id")
	if er != nil {
		return er
	}
	return a.UserService.Disable(c.Context(), int64(id))
}

func (a *Sapi) UserDel(c *fiber.Ctx) error {
	id, er := c.ParamsInt("id")
	if er != nil {
		return er
	}
	return a.UserService.Delete(c.Context(), int64(id))
}

func (a *Sapi) UserList(c *fiber.Ctx) error {
	var opt types.UserOption
	er := c.QueryParser(&opt)
	if er != nil {
		return er
	}
	users, er := a.UserService.List(c.Context(), &opt)
	if er != nil {
		return er
	}

	return c.JSON(users)
}

func (a *Sapi) TaskDo(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	uid := c.Locals("userId").(int64)
	return a.TaskService.DoOnce(c.Context(), uid, int64(id))
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
	uid := c.Locals("userId").(int64)

	return a.TaskService.Update(c.Context(), uid, &t)
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
	uid := c.Locals("userId").(int64)
	return a.TaskService.Start(c.Context(), uid, int64(id))
}

func (a *Sapi) TaskStop(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	uid := c.Locals("userId").(int64)
	return a.TaskService.Stop(c.Context(), uid, int64(id))
}

func (a *Sapi) TaskDel(c *fiber.Ctx) error {
	id, _ := c.ParamsInt("id")
	uid := c.Locals("userId").(int64)
	return a.TaskService.Del(c.Context(), uid, int64(id))
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
	uid := c.Locals("userId").(int64)
	count, er := a.TaskService.StartBatch(c.Context(), uid, &tbatch)
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
	uid := c.Locals("userId").(int64)
	count, er := a.TaskService.StopBatch(c.Context(), uid, &tbatch)
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
	uid := c.Locals("userId").(int64)
	count, er := a.TaskService.DelBatch(c.Context(), uid, &tbatch)
	if er != nil {
		return er
	}
	return c.JSON(fiber.Map{
		"count": count,
	})
}

func (a *Sapi) ExcutePage(c *fiber.Ctx) error {
	var opt types.TaskExcuteOption

	er := c.QueryParser(&opt)
	if er != nil {
		return er
	}
	page, er := a.TaskService.ExcutePageByOption(c.Context(), &opt)
	if er != nil {
		return er
	}
	return c.JSON(page)
}

func (a *Sapi) Dashboard(c *fiber.Ctx) error {
	data := a.TaskService.Dashboard(c.Context())

	return c.JSON(data)
}
