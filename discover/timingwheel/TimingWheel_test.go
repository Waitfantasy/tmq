package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

func TestNewTimingWheel(t *testing.T) {
	tw, _ := NewTimingWheel(time.Millisecond, 200)
	fmt.Println(tw.interval, tw.tick, tw.currentTime)
}
