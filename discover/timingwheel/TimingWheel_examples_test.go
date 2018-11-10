package timingwheel

import (
	"fmt"
	"time"
)

func ExampleNewTimingWheel() {
	tw, _ := NewTimingWheel(time.Millisecond, 200)
	fmt.Println(tw)
}
