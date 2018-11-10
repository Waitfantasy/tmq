package timingwheel

import (
	"sync"
	"time"
)

type TimingWheel struct {
	wg                  sync.WaitGroup
	tick                uint64
	wheels              uint64
	interval            uint64
	currentTime         uint64
	slots               []*slot
	dq                  *delayQueue
	overflowTimingWheel *TimingWheel
}

func NewTimingWheel(tick time.Duration, wheels uint64) (*TimingWheel, error) {
	tickMs := tick / time.Millisecond
	if tickMs < 0 {
		// 默认tick为1毫秒
		tickMs = 1
		//return nil, errors.New("tick must >= 1 million second")
	}

	return newTimingWheel(uint64(tickMs),
		wheels,
		truncate(ms(time.Now()), uint64(tickMs)),
		newDelayQueue(int(wheels))), nil

}

func newTimingWheel(tick, wheels, curTime uint64, dq *delayQueue) *TimingWheel {
	tw := &TimingWheel{}
	tw.wheels = wheels
	tw.tick = tick
	tw.interval = tw.tick * tw.wheels
	tw.currentTime = curTime

	// 初始化buckets
	tw.slots = make([]*slot, wheels)
	for i := 0; i < int(wheels); i++ {
		tw.slots[i] = newBucket()
	}

	// 初始化delay queue
	tw.dq = dq

	return tw
}

func (tw *TimingWheel) add(task *Task) bool {
	// 过期任务
	if task.expiration < tw.currentTime+tw.tick {
		return false
	}

	// 当前时间轮可以插入该任务
	if task.expiration < tw.currentTime+tw.interval {
		// 找出可插入任务的slots
		taskTick := task.expiration / tw.tick
		slot := tw.slots[taskTick%tw.wheels]
		slot.add(task)

		if slot.setExpiration(task.expiration) {
			tw.dq.Offer(slot)
		}

		return true
	}

	// 当前时间轮承载不了任务到期时间
	tw.overflowTimingWheel = newTimingWheel(tw.interval, tw.wheels, tw.currentTime, tw.dq)

	return true
}

func (tw *TimingWheel) Start() {
	// 运行线程对任务进行轮询
	go func() {
		tw.dq.Poll()
	}()

	// 运行线程接收可执行的任务
	go func() {
		slot := <- tw.dq.ExpireChan
		// 推进时钟
		// 执行任务调度
	}()
}
