package user

import (
	"context"
	"github.com/eininst/scheduler/internal/model"
	"github.com/eininst/scheduler/internal/service"
	"github.com/eininst/scheduler/internal/util"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)

type UserService interface {
	Login(ctx context.Context, username string, password string) (*model.SchedulerUser, error)
	List(ctx context.Context) ([]*model.SchedulerUser, error)
}

type userService struct {
	Store *session.Store
	DB    *gorm.DB      `inject:""`
	Rcli  *redis.Client `inject:""`
}

func NewService() UserService {
	return &userService{}
}

func (us *userService) Login(ctx context.Context,
	username string, password string) (*model.SchedulerUser, error) {

	var u model.SchedulerUser
	us.DB.First(&u, "name = ?", username)
	if u.Id == 0 {
		return nil, service.NewServiceError("账号或密码错误")
	}
	if u.Password != util.Md5(password) {
		return nil, service.NewServiceError("账号或密码错误")
	}
	return &u, nil
}

func (us *userService) Add(ctx context.Context, user *model.SchedulerUser) error {
	return nil
}

func (us *userService) List(ctx context.Context) ([]*model.SchedulerUser, error) {
	var users []*model.SchedulerUser
	us.DB.WithContext(ctx).Find(&users, "status = 'ok'")
	return users, nil
}
