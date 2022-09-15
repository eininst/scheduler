package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
)

func main() {
	parse := cron.NewParser(
		cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)
	s, er := parse.Parse("* * * * * *")

	fmt.Println(s)
	fmt.Println(er)
	//c := cron.New(cron.WithSeconds())
	////
	////fmt.Println(123)
	////
	//_, err := c.AddFunc("*/* 1 * * * *", func() {
	//	fmt.Println(1)
	//})
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//c.Start()
	//
	//<-make(chan int)
}
