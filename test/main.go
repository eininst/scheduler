package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	c := cron.New(cron.WithSeconds())

	fmt.Println(123)
	c.AddFunc("*/* 1 * * * *", func() {
		fmt.Println(1)
	})
	c.Start()

	<-make(chan int)
}
