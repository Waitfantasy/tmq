package timingwheel

import (
	"container/list"
	"sync/atomic"
	"unsafe"
)

type TaskEntry struct {
	f          func()
	el         *list.Element
	slot       unsafe.Pointer
	cancelled  bool
	expiration int64
}

func (t *TaskEntry) Stop() bool {
	for b := t.getBucket(); b != nil; b = t.getBucket() {
		t.cancelled = b.remove(t)
	}
	return t.cancelled
}

func (t *TaskEntry) setBucket(b *bucket) {
	atomic.StorePointer(&t.slot, unsafe.Pointer(b))
}

func (t *TaskEntry) getBucket() *bucket {
	return (*bucket)(atomic.LoadPointer(&t.slot))
}
