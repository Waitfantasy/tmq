package timingwheel

import "time"

func ms(t time.Time) uint64 {
	return uint64(time.Duration(t.UnixNano()) / time.Millisecond)
}

func truncate(x, m uint64) uint64 {
	if m <= 0 {
		return x
	}

	return x - x%m
}
