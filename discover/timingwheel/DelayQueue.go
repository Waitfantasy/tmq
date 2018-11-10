package timingwheel

import (
	"container/heap"
	"sync"
	"time"
)

type item struct {
	Slot   interface{}
	Expire uint64
	Index  int
}

type timingWheelHeap []*item

func newTimingWheelHeap(size int) timingWheelHeap {
	return make(timingWheelHeap, 0, size)
}

func (pq timingWheelHeap) Len() int {
	return len(pq)
}

func (pq timingWheelHeap) Less(i, j int) bool {
	return pq[i].Expire < pq[j].Expire
}

func (pq timingWheelHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *timingWheelHeap) Push(x interface{}) {
	if item, ok := x.(*item); ok {
		item.Index = len(*pq)
		*pq = append(*pq, item)
	}
	//l := len(*pq)
	//c := cap(*pq)
	//// 即将超过容量
	//if l + 1 > c {
	//	// 分配2倍容量
	//	v := make(timingWheelHeap, l, c * 2)
	//	copy(v, *pq)
	//	*pq = v
	//}
	//
	//*pq = (*pq)[0 : l+1]
	//if item, ok := x.(*item); ok {
	//	item.Index = l
	//	(*pq)[l] = item
	//}
}

func (pq *timingWheelHeap) Pop() interface{} {
	tmp := *pq
	l := len(tmp)
	item := tmp[l-1]
	item.Index = -1
	*pq = tmp[0 : l-1]
	return item
}

type delayQueue struct {
	mu        sync.Mutex
	pq        *timingWheelHeap
	available chan struct{}
}

func newDelayQueue(size int) *delayQueue {
	h := newTimingWheelHeap(size)
	return &delayQueue{
		pq:        &h,
		available: make(chan struct{}),
	}
}

func (dq *delayQueue) offer(s *slot) {
	item := &item{Slot: s, Expire: s.expiration}
	dq.mu.Lock()
	heap.Push(dq.pq, item)
	dq.mu.Unlock()

	// 如果是队首元素, 通知线程
	if (*dq.pq)[0] == item {

	}
}

// now 表示从当前时间开始等待
func (dq *delayQueue) peekAndRemove(now uint64) (*item, uint64) {
	if dq.pq.Len() == 0 {
		return nil, 0
	}

	first := (*dq.pq)[0]
	// 超期
	if first.Expire > now {
		return nil, first.Expire - now
	}

	first = heap.Pop(dq.pq).(*item)
	return first, 0
}

func (dq *delayQueue) Poll() {
	for {
		first, delay := dq.peekAndRemove(ms(time.Now()))

		if first == nil {
			if delay == 0 {
				select {

				}
 			} else if delay > 0{

			}
		} else {

		}

		panic("impl me")
	}
}
