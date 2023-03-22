package main

import (
	"fmt"
	"time"

	"github.com/apigear-io/scheduler"
)

func main() {
	s := scheduler.New(time.Millisecond * 10)
	s.CreateJob().Repeat(10).Every(time.Millisecond * 100).Do(func(t, dt int64) {
		fmt.Printf("tick: %d, delta-tick: %d\n", t, dt)
	})
	s.OnJobRemove(func(j *scheduler.Job) {
		fmt.Printf("job removed: %s\n", j.ID())
		go s.Stop()
	})
	s.Run()
}
