package gojm

import (
	"context"
	"sync"
	"time"
)

// Job is the wrapper of a function which allows to wait the result
// asynchronously.
type Job struct {
	mutex sync.Mutex
	id    any
	f     func(ctx context.Context) *JobResult

	result  *JobResult
	resultC chan *JobResult
}

// NewJob wraps a function to Job.
func NewJob(f func(ctx context.Context) *JobResult) *Job {
	return &Job{
		id:      nil,
		mutex:   sync.Mutex{},
		f:       f,
		resultC: make(chan *JobResult, 1),
	}
}

// IsCompleted returns true if the job already completed.
func (j *Job) IsCompleted() bool {
	return j.GetResult() != nil
}

// GetResult returns the JobResult if the job already completed. Otherwise, it
// returns nil.
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

// WaitResult returns the result if the job already completed. Otherwise, it
// blocks the current process until the job completes.
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

// Exec executes the job and returns the result. Do not call this method
// multiple times.
func (j *Job) Exec(ctx context.Context) *JobResult {
	result := j.f(ctx)
	if result == nil {
		result = EmptyResult()
	}

	result.DoneAt = time.Now()

	j.resultC <- result
	close(j.resultC)

	return result
}
