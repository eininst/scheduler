package main

import (
	"github.com/eininst/scheduler"
	"os"
)

func main() {
	config := os.Getenv("config")

	if config == "" {
		config = "/config.yaml"
	}

	app := scheduler.New(config)
	app.Listen()
}
