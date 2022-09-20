package main

import (
	"fmt"
	"github.com/ivpusic/grpool"
	"time"
)

func test() {
	time.Sleep(time.Second * 3)
	fmt.Println(123)
}

func w(pool *grpool.Pool) {
	go func() {

		pool.JobQueue <- func() {
			test()
		}
	}()
}
func main() {
	pool := grpool.NewPool(1, 0)
	for i := 0; i < 10; i++ {
		w(pool)
	}

	<-make(chan int)
}
