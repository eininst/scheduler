package conf

import (
	"github.com/eininst/go-jwt"
	"github.com/eininst/ninja"
	"github.com/eininst/rlock"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/data"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
)

func Inject() {
	//inject resources
	rcli := data.NewRedisClient()
	ninja.Provide(rcli)
	ninja.Provide(rlock.New(rcli))

	db := data.NewDB()
	ninja.Provide(db)

	ninja.Provide(jwt.New(configs.Get("secretKey").String()))

	//inject services
	ninja.Provide(user.NewService())
	ninja.Provide(task.NewService())
}
