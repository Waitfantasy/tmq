package timingwheel

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

type TimingWheel struct {
	tick                int64
	wheels              int64
	interval            int64
	currentTime         int64
	buckets             []*bucket
	dq                  *delayQueue
	wg                  sync.WaitGroup
	exitChan            chan struct{}
	overflowTimingWheel unsafe.Pointer
}

func NewTimingWheel(tick time.Duration, wheels int64) *TimingWheel {
	tickMs := tick / time.Millisecond
	if tickMs < 0 {
		// 默认tick为1毫秒
		tickMs = 1
		//return nil, errors.New("tick must >= 1 million second")
	}

	return newTimingWheel(
		int64(tickMs),
		wheels,
		truncate(ms(time.Now()), int64(tickMs)),
		newDelayQueue(int(wheels)),
	)
}

func newTimingWheel(tick, wheels, curTime int64, dq *delayQueue) *TimingWheel {
	tw := &TimingWheel{
		tick:        tick,
		wheels:      wheels,
		interval:    tick * wheels,
		currentTime: curTime,
		dq:          dq,
		exitChan:    make(chan struct{}),
	}

	// 初始化buckets
	tw.buckets = make([]*bucket, wheels)
	for i := 0; i < int(wheels); i++ {
		tw.buckets[i] = newBucket()
	}

	return tw
}

func (tw *TimingWheel) add(task *TaskEntry) bool {
	// 过期任务
	if task.expiration < tw.currentTime+tw.tick {
		return false
	}

	// 当前时间轮可以插入该任务
	if task.expiration < tw.currentTime+tw.interval {
		// 找出可插入任务的slots
		taskTick := task.expiration / tw.tick

		b := tw.buckets[taskTick%tw.wheels]
		b.add(task)
		if b.setExpiration(task.expiration) {
			tw.dq.Offer(b)
		}

		return true
	}

	// 当前任务超过时间轮周期, 扩大时间轮的tick
	otw := atomic.LoadPointer(&tw.overflowTimingWheel)
	if otw == nil {
		atomic.CompareAndSwapPointer(
			&tw.overflowTimingWheel,
			nil,
			unsafe.Pointer(newTimingWheel(
				tw.interval,
				tw.wheels,
				tw.currentTime,
				tw.dq,
			)),
		)
		otw = atomic.LoadPointer(&tw.overflowTimingWheel)
	}

	return (*TimingWheel)(otw).add(task)
}

// 推进时钟前进
func (tw *TimingWheel) advanceClock(expiration int64) {
	if expiration >= tw.currentTime+tw.tick {
		tw.currentTime = truncate(expiration, tw.tick)
	}

	if otw := atomic.LoadPointer(&tw.overflowTimingWheel); otw != nil {
		(*TimingWheel)(otw).advanceClock(expiration)
	}
}

func (tw *TimingWheel) Start() {
	// 运行线程对任务进行分发
	tw.wg.Add(2)
	go func() {
		tw.dq.Dispatch(tw.exitChan)
		tw.wg.Done()

	}()

	// 运行线程接收可执行的任务
	go func() {
		for {
			select {
			case slot := <-tw.dq.C:
				tw.advanceClock(slot.expiration)
				slot.flush(tw.add)
			case <-tw.exitChan:
				tw.wg.Done()
				return
			}
		}
	}()

}

func (tw *TimingWheel) Stop() {
	close(tw.exitChan)
	tw.wg.Wait()
}

func (tw *TimingWheel) AfterFunc(d time.Duration, f func()) *TaskEntry {
	t := &TaskEntry{
		f:          f,
		expiration: ms(time.Now().Add(d)),
	}

	if !tw.add(t) {
		go t.f()
	}

	return t
}
