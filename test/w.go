package main

import (
	"fmt"
	"github.com/eininst/rlock"
	"github.com/eininst/scheduler/internal/data"
)

type Wk struct {
}

func (w Wk) A() {
	fmt.Println(1)
}
func main() {
	rcli := data.NewRedisClient()
	rlock.SetDefault(rcli)
	//ninja.Provide(rcli)
	//ninja.Provide(rlock.New(rcli))
	//ninja.Provide(data.NewRsClient(rcli))

	//ok, cancel := rlock.TryAcquire(lockName, time.Second*time.Duration(1))
	//flog.Info(ok)
	//flog.Info(cancel)
}
