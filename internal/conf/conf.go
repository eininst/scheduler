package conf

import (
	"github.com/eininst/ninja"
	"github.com/eininst/rlock"
	"github.com/eininst/scheduler/internal/data"
)

func Inject() {
	//inject resources
	rcli := data.NewRedisClient()
	ninja.Provide(rcli)
	ninja.Provide(rlock.New(rcli))

	db := data.NewDB()
	ninja.Provide(db)
}
