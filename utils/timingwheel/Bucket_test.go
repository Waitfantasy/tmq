package timingwheel

import (
	"sync"
	"testing"
	"time"
)

func TestBucket_Add(t *testing.T) {
	durations := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
	}

	b := newBucket()
	wg := sync.WaitGroup{}
	for i := 0; i < len(durations); i++ {
		wg.Add(1)
		go func(i int) {
			t := &TaskEntry{
				expiration: int64(durations[i]),
			}

			b.Add(t)

			wg.Done()
		}(i)
	}

	wg.Wait()
	if b.tasks.Len() != len(durations) {
		t.Error()
		return
	}
}

func TestBucket_Remove(t *testing.T) {
	durations := []time.Duration{
		1 * time.Millisecond,
		2 * time.Millisecond,
		3 * time.Millisecond,
		4 * time.Millisecond,
	}

	b := newBucket()
	wg := sync.WaitGroup{}
	for i := 0; i < len(durations); i++ {
		wg.Add(1)
		go func(i int) {
			t := &TaskEntry{
				expiration: int64(durations[i]),
			}

			b.Add(t)

			wg.Done()
		}(i)
	}

	wg.Wait()
	if b.tasks.Len() != len(durations) {
		t.Error()
		return
	}

	wg.Add(2)
	go func() {
		e := b.tasks.Front()
		for e != nil {
			next := e.Next()
			b.remove(e.Value.(*TaskEntry))
			e = next
		}

		wg.Done()
	}()

	go func() {
		e := b.tasks.Front()
		for e != nil {
			next := e.Next()
			b.remove(e.Value.(*TaskEntry))
			e = next
		}
		wg.Done()
	}()

	wg.Wait()
	if b.tasks.Len() != 0 {
		t.Error()
		return
	}
}