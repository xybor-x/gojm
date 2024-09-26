package gojm

import "time"

type Priority struct {
	value          int
	agingTimeSlice time.Duration
}

func NewPriority(p int) Priority {
	return Priority{value: p, agingTimeSlice: 0}
}

func (p Priority) WithAging(timeslice time.Duration) Priority {
	p.agingTimeSlice = timeslice
	return p
}
