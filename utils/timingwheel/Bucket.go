package timingwheel

import (
	"container/list"
	"sync"
	"sync/atomic"
)

type bucket struct {
	mu         sync.Mutex
	expiration int64
	tasks      *list.List
}

func newBucket() *bucket {
	b := new(bucket)
	b.tasks = list.New()
	return b
}

func (b *bucket) Add(t *TaskEntry) {
	b.add(t)
}

func (b *bucket) add(t *TaskEntry) {
	b.mu.Lock()
	t.el = b.tasks.PushBack(t)
	b.mu.Unlock()
	t.setBucket(b)
}

func (b *bucket) Remove(t *TaskEntry) bool {
	return b.remove(t)
}

func (b *bucket) remove(t *TaskEntry) bool {
	if t.getBucket() != b {
		return false
	}

	b.mu.Lock()
	b.tasks.Remove(t.el)
	t.el = nil
	b.mu.Unlock()
	t.setBucket(nil)
	return true
}

func (b *bucket) setExpiration(expiration int64) bool {
	if atomic.SwapInt64(&b.expiration, expiration) == expiration {
		return false
	}
	return true
}

func (b *bucket) flush(reinsert func(entry *TaskEntry) bool) {
	for head := b.tasks.Front(); head != nil; {
		next := head.Next()
		t := head.Value.(*TaskEntry)
		b.remove(t)
		if !reinsert(t) && !t.cancelled {
			go t.f()
		}
		head = next
	}
}
