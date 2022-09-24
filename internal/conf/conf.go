package conf

import (
	"github.com/eininst/go-jwt"
	"github.com/eininst/ninja"
	"github.com/eininst/rlock"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/data"
	"github.com/eininst/scheduler/internal/service/mail"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/robfig/cron/v3"
)

func Inject() {
	//inject resources
	rcli := data.NewRedisClient()
	ninja.Provide(rcli)
	ninja.Provide(rlock.New(rcli))
	ninja.Provide(data.NewRsClient(rcli))

	db := data.NewDB()
	ninja.Provide(db)

	ninja.Provide(jwt.New(configs.Get("secretKey").String()))

	cronCli := cron.New(cron.WithSeconds())
	ninja.Provide(cronCli)

	//inject services
	ninja.Provide(user.NewService())
	ninja.Provide(task.NewService())

	ninja.Provide(mail.NewService())
}
