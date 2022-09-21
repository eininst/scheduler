package main

import (
	"github.com/eininst/flog"
	"time"
)

func main() {
	f := time.Now().Add(-time.Hour * 24 * 3).Format("2006-01-02 15:04:05")
	flog.Info(f)
}
