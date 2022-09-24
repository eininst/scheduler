package main

import (
	"github.com/eininst/scheduler"
)

func main() {
	//cfgPath := flag.String("cfg", "", "a config path")
	//flag.Parse()
	//
	//if *cfgPath == "" {
	//	log.Fatal("require set a config, eg: -cfg=xxxx")
	//}
	//app := scheduler.New(*cfgPath)
	app := scheduler.New("./configs/config.yml")
	app.Listen(":3000")
}
