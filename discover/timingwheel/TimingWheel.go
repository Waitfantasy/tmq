package timingwheel

import (
	"time"
)

type TimingWheel struct {
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
	for i := 0; i < int(wheels); i ++ {
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
		// 找出可防止任务的slots
		taskTick := task.expiration / tw.tick
		slot := tw.slots[taskTick%tw.wheels]
		slot.add(task)


		if slot.setExpiration(task.expiration) {
			//
		}
	}

	// 当前时间轮承载不了任务到期时间
	tw.overflowTimingWheel = newTimingWheel(tw.interval, tw.wheels, tw.currentTime, tw.dq)

	return true
}

func (tw *TimingWheel) newOverflowTimingWheel(tick, wheels uint64) *TimingWheel {
	otw := new(TimingWheel)
	otw.tick = tick
	otw.wheels = wheels
	otw.interval = otw.wheels * otw.tick
	return otw
}
