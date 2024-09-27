package gojm

import (
	"context"
	"sync"
	"time"

	priorityqueue "github.com/xybor-x/priority_queue"
)

// JobManager allows to schedule job based on its priority.
type JobManager struct {
	hook  func(ctx context.Context, job JobWrapper)
	queue *priorityqueue.PriorityQueue[*Job]
}

// New initialized a JobManager.
func New() *JobManager {
	return &JobManager{
		hook:  nil,
		queue: priorityqueue.Default[*Job](),
	}
}

// AddPriority sets a new Priority to JobManager.
func (m *JobManager) AddPriority(p Priority) {
	if err := m.queue.SetPriority(p, p.value); err != nil {
		panic(err)
	}

	if p.agingTimeSlice != nil {
		if err := m.queue.SetAgingTimeSlice(p, *p.agingTimeSlice); err != nil {
			panic(err)
		}
	}
}

// SetCommongJobAging sets a timeslice. When the job has existed for more than
// this timeslice, it will be moved to the higher priority. This timeslice is
// only applied when the priority hasn't its own aging.
func (m *JobManager) SetCommonJobAging(timeslice time.Duration) {
	if err := m.queue.SetCommonAgingTimeSlice(timeslice); err != nil {
		panic(err)
	}
}

// RefreshEvery sets an interval which refreshes the priority of jobs every time
// interval passed. If you do not call this method, this interval will be chosen
// automatically (equal to the least aging timeslice of all priorities).
func (m *JobManager) RefreshEvery(interval time.Duration) {
	m.queue.SetAgingInterval(interval)
}

// Schedule adds a Job to JobManager, the job will be scheduled to execute
// later.
func (m *JobManager) Schedule(priority Priority, job *Job) error {
	return m.queue.Enqueue(priority, job)
}

// Run starts the JobManager. The parameter numThreads specifies the number of
// Jobs which could be executed concurrently.
func (m *JobManager) Run(ctx context.Context, numThreads int) error {
	wg := sync.WaitGroup{}
	wg.Add(numThreads)

	var finalErr error
	var mutex sync.Mutex

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for i := 0; i < numThreads; i++ {
		go func() {
			defer wg.Done()

			if err := m.RunOne(ctx); err != nil {
				mutex.Lock()

				if finalErr == nil {
					finalErr = err
					cancel()
				}

				mutex.Unlock()
			}
		}()
	}

	wg.Wait()
	return finalErr
}

// RunOne starts the JobManager which only one Job could be executed at a time.
func (m *JobManager) RunOne(ctx context.Context) error {
	for {
		job, err := m.queue.WaitDequeue(ctx)
		if err != nil {
			return err
		}

		job.To().Exec(ctx)
		if m.hook != nil {
			m.hook(ctx, wrapJob(job))
		}
	}
}

// Hook sets a trigger function to be executed when the job has just completed.
// This method should be called before calling Run() or RunOne().
func (m *JobManager) Hook(trigger func(ctx context.Context, job JobWrapper)) {
	m.hook = trigger
}
