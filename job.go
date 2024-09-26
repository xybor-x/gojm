package gojm

import (
	"context"
	"sync"
	"time"
)

type Job struct {
	mutex sync.Mutex
	id    any
	f     func(ctx context.Context) *JobResult

	result  *JobResult
	resultC chan *JobResult
}

func NewJob(f func(ctx context.Context) *JobResult) *Job {
	return &Job{
		id:      nil,
		mutex:   sync.Mutex{},
		f:       f,
		resultC: make(chan *JobResult, 1),
	}
}

func (j *Job) GetResult() *JobResult {
	if j.result != nil {
		return j.result
	}

	select {
	case result := <-j.resultC:
		j.mutex.Lock()
		if j.result == nil {
			j.result = result
		}
		j.mutex.Unlock()
	default:
	}

	return j.result
}

func (j *Job) WaitResult(ctx context.Context) *JobResult {
	if result := j.GetResult(); result != nil {
		return result
	}

	select {
	case result := <-j.resultC:
		j.mutex.Lock()
		if j.result == nil {
			j.result = result
		}
		j.mutex.Unlock()

		return j.GetResult()
	case <-ctx.Done():
		return nil
	}
}

func (j *Job) Execute(ctx context.Context) *JobResult {
	result := j.f(ctx)
	result.DoneAt = time.Now()

	j.resultC <- result
	close(j.resultC)

	return result
}
