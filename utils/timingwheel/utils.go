package timingwheel

import "time"

func ms(t time.Time) int64 {
	return int64(time.Duration(t.UnixNano()) / time.Millisecond)
}

func truncate(x, m int64) int64 {
	if m <= 0 {
		return x
	}

	return x - x%m
}
