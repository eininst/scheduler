package conf

import (
	"github.com/eininst/go-jwt"
	"github.com/eininst/rlock"
	"github.com/eininst/scheduler/configs"
	"github.com/eininst/scheduler/internal/data"
	"github.com/eininst/scheduler/internal/inject"
	"github.com/eininst/scheduler/internal/service/mail"
	"github.com/eininst/scheduler/internal/service/task"
	"github.com/eininst/scheduler/internal/service/user"
	"github.com/robfig/cron/v3"
)

func Inject() {
	//inject resources
	rcli := data.NewRedisClient()
	inject.Provide(rcli)
	inject.Provide(rlock.New(rcli))
	inject.Provide(data.NewRsClient(rcli))

	db := data.NewDB()
	inject.Provide(db)

	inject.Provide(jwt.New(configs.Get("secretKey").String()))

	cronCli := cron.New(cron.WithSeconds())
	inject.Provide(cronCli)

	//inject services
	inject.Provide(user.NewService())
	inject.Provide(task.NewService())

	inject.Provide(mail.NewService())
}
