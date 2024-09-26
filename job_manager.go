package gojm

import (
	"context"
	"sync"
	"time"

	priorityqueue "github.com/xybor-x/priority_queue"
)

type JobManager struct {
	hook  func(ctx context.Context, job priorityqueue.Element[*Job])
	queue *priorityqueue.PriorityQueue[*Job]
}

func New() *JobManager {
	return &JobManager{
		hook:  nil,
		queue: priorityqueue.Default[*Job](),
	}
}

func (m *JobManager) AddPriority(p Priority) {
	if err := m.queue.SetPriority(p, p.value); err != nil {
		panic(err)
	}

	if err := m.queue.SetAgingTimeSlice(p, p.agingTimeSlice); err != nil {
		panic(err)
	}
}

func (m *JobManager) SetCommonJobAging(timeslice time.Duration) {
	if err := m.queue.SetCommonAgingTimeSlice(timeslice); err != nil {
		panic(err)
	}
}

func (m *JobManager) RefreshEvery(interval time.Duration) {
	m.queue.SetAgingInterval(interval)
}

func (m *JobManager) Schedule(priority Priority, job *Job) error {
	return m.queue.Enqueue(priority, job)
}

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

func (m *JobManager) RunOne(ctx context.Context) error {
	for {
		job, err := m.queue.WaitDequeue(ctx)
		if err != nil {
			return err
		}

		job.To().Execute(ctx)
		if m.hook != nil {
			m.hook(ctx, job)
		}
	}
}

func (m *JobManager) Hook(trigger func(ctx context.Context, job priorityqueue.Element[*Job])) {
	m.hook = trigger
}
