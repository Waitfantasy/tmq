package timingwheel

import (
	"container/heap"
	"fmt"
	"math/rand"
	"testing"
)

func TestTimingWheelHeap(t *testing.T) {
	h := newTimingWheelHeap(10)
	ph := &h
	heap.Init(ph)
	for i := 0; i < 100; i++ {
		r := rand.Intn(1000)
		heap.Push(ph, &item{
			Slot:   r,
			Expire: uint64(r),
		})
	}

	for i:= 0; i< 100; i++ {
		v := heap.Pop(ph)
		fmt.Println(v.(*item).Slot.(int))
	}
}