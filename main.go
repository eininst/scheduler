package main

import (
	"fmt"
	"github.com/eininst/flog"
	"github.com/eininst/rlock"
	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"runtime"
	"time"
)

func GetRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:         "localhost:6379",
		DB:           0,
		DialTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     100,
		MinIdleConns: 25,
		PoolTimeout:  30 * time.Second,
	})
}

func init() {
	rlock.SetDefault(GetRedis())
}
func main() {

	//11.123 -> 1s -> 12.123
	//11.500 -> 1s -> 12.500

	flog.SetTimeFormat("2006/01/02 15:04:05.000000")
	//
	fmt.Println(runtime.NumCPU())
	c := cron.New(cron.WithSeconds())
	c.Start()

	time.Sleep(time.Second)
	go func() {
		eid, er := c.AddFunc("*/5 * * * * *", func() {
			flog.Info("My name is 1")
		})
		flog.Info(eid, er)
		time.Sleep(time.Millisecond * 200)
		c.AddFunc("*/5 * * * * *", func() {
			flog.Info("My name is 2")
		})

		time.Sleep(time.Millisecond * 567)
		c.AddFunc("*/5 * * * * *", func() {
			flog.Info("My name is 5")
		})
	}()
	<-make(chan int)
	ctx := c.Stop()
	<-ctx.Done()
}
