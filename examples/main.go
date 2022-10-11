package main

import (
	"flag"
	"github.com/eininst/scheduler"
	"log"
)

func main() {
	cfgPath := flag.String("cfg", "", "a config path")
	flag.Parse()

	if *cfgPath == "" {
		log.Fatal("require set a config, eg: -cfg=xxxx")
	}

	app := scheduler.New(*cfgPath)

	app.Listen()
}
