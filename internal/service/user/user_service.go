package user

import (
	"context"
	"github.com/eininst/scheduler/internal/service"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserService interface {
}

type userService struct {
	Db   *gorm.DB      `inject:""`
	Rcli *redis.Client `inject:""`
}

func (t *userService) Add(ctx context.Context, user *service.SchedulerUser) error {
	return nil
}
