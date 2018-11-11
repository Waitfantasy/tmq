package timingwheel

import (
	"container/heap"
	"sync"
	"time"
)

type item struct {
	Slot   interface{}
	Expire int64
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
}

func (pq *timingWheelHeap) Push(x interface{}) {
	if item, ok := x.(*item); ok {
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
	*pq = tmp[0 : l-1]
	return item
}

func (pq *timingWheelHeap) Peek() *item {
	if len(*pq) == 0 {
		return nil
	}

	return (*pq)[0]
}

type delayQueue struct {
	C         chan *bucket
	mu        sync.RWMutex
	pq        *timingWheelHeap
	available chan struct{}
	ExitChan  chan struct{}
}

func newDelayQueue(size int) *delayQueue {
	h := newTimingWheelHeap(size)
	return &delayQueue{
		C:         make(chan *bucket),
		pq:        &h,
		available: make(chan struct{}),
		ExitChan:  make(chan struct{}),
	}
}

func (dq *delayQueue) Offer(s *bucket) {
	item := &item{Slot: s, Expire: s.expiration}
	dq.mu.Lock()
	heap.Push(dq.pq, item)
	dq.mu.Unlock()

	// 如果是队首元素, 通知线程
	dq.mu.RLock()
	first := (*dq.pq)[0]
	dq.mu.RUnlock()
	if first == item {
		go func() {
			dq.available <- struct{}{}
			return
		}()
	}
}

// 通过时间来驱动poll
func (dq *delayQueue) Poll(timestamp int64, unit time.Duration) *item {
	waitTime := time.Duration(timestamp) * unit
	for {
		dq.mu.RLock()
		first := dq.pq.Peek()
		dq.mu.RUnlock()
		if first == nil {
			if waitTime <= 0 {
				return nil
			} else {
				select {
				case <-time.After(waitTime):
					continue
				}
			}
		} else {
			delay := int(first.Expire - ms(time.Now()))
			if delay <= 0 {
				return heap.Pop(dq.pq).(*item)
			}

			first = nil
			if waitTime < time.Duration(delay) {
				select {
				case <-time.After(waitTime):
					continue
				}
			} else {
				select {
				case <-time.After(time.Duration(delay) * time.Millisecond):
					continue
				}
			}
		}
	}
}

func (dq *delayQueue) Dispatch(exit chan struct{}) {
	for {
		dq.mu.RLock()
		first := dq.pq.Peek()
		dq.mu.RUnlock()
		// 元素为空
		if first == nil {
			select {
			case <-dq.available:
				continue
			case <-exit:
				goto EXIT
			}
		} else {
			delay := first.Expire - ms(time.Now())
			first = nil
			// 还未到期
			if delay > 0 {
				select {
				case <-dq.available:
					continue
				case <-time.After(time.Duration(delay) * time.Millisecond): // 等待
					continue
				case <-exit:
					goto EXIT
				}
			} else {
				dq.C <- heap.Pop(dq.pq).(*item).Slot.(*bucket)
			}
		}
	}
EXIT:
	return
}
