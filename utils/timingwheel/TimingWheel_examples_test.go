package timingwheel

import (
	"fmt"
	"sync"
	"time"
)

func Example_AfterFunc() {
	wg := sync.WaitGroup{}
	tw := NewTimingWheel(time.Millisecond, 20)
	tw.Start()
	defer tw.Stop()
	durations := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
		5 * time.Millisecond,
	}
	expires := make(chan time.Duration, len(durations))
	start := time.Now()
	for i := 0; i < len(durations); i++ {
		wg.Add(1)
		tw.AfterFunc(durations[i], func() {
			expires <- time.Since(start)
			wg.Done()
		})
	}

	go func() {
		for expire := range expires {
			fmt.Println(expire)
		}
	}()

	wg.Wait()
	close(expires)

}
