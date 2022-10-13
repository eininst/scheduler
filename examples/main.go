package main

import "github.com/eininst/scheduler"

func main() {
	app := scheduler.New("./configs/config.yaml")
	app.Listen()
}
