package timingwheel

import (
	"container/list"
	"sync/atomic"
)

type Task struct {
	f          func()
	el         *list.Element
	expiration uint64
}

type slot struct {
	expiration uint64
	tasks      *list.List
}

func newBucket() *slot {
	b := new(slot)
	b.tasks = list.New()
	return b
}

func (b *slot) add(t *Task) {
	b.tasks.PushBack(t)
}

func (b *slot) setExpiration(expiration uint64) bool {
	if atomic.SwapUint64(&b.expiration, expiration) == expiration {
		return false
	}
	return true
}
