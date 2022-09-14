package task

import (
	"context"
	"github.com/eininst/scheduler/internal/model"
	"github.com/go-redis/redis/v8"
)

type TaskService interface {
}
type taskService struct {
	Rcli *redis.Client `inject:""`
}

func (t *taskService) Add(ctx context.Context, task *model.SchedulerTask) error {
	return nil
}
