package task

import (
	"context"
	"github.com/eininst/scheduler/internal/types"
	"github.com/go-redis/redis/v8"
)

type TaskService interface {
}
type taskService struct {
	Rcli *redis.Client `inject:""`
}

func (t *taskService) Add(ctx context.Context, task *types.Task) error {
	return nil
}
