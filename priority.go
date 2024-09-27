package gojm

import (
	"fmt"
	"time"
)

type Priority struct {
	name           string
	value          int
	agingTimeSlice *time.Duration
}

func NewPriority(name string, p int) Priority {
	return Priority{name: name, value: p, agingTimeSlice: nil}
}

func (p Priority) WithAging(timeslice time.Duration) Priority {
	p.agingTimeSlice = &timeslice
	return p
}

func (p Priority) WithNoAging() Priority {
	zeroDuration := time.Duration(0)
	p.agingTimeSlice = &zeroDuration
	return p
}

func (p Priority) Name() string {
	return p.name
}

func (p Priority) Value() int {
	return p.value
}

func (p Priority) String() string {
	return fmt.Sprintf("%s(%d)", p.Name(), p.Value())
}
