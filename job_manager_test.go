package gojm_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/gojm"
	priorityqueue "github.com/xybor-x/priority_queue"
)

var Urgent = gojm.NewPriority(0)
var Necessary = gojm.NewPriority(10).WithAging(1 * time.Second)
var Background = gojm.NewPriority(100)

func Test_JobManager_RunOne(t *testing.T) {
	jm := gojm.New()
	jm.AddPriority(Urgent)
	jm.AddPriority(Necessary)
	jm.AddPriority(Background)

	urgentJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.EmptyResult()
	})

	neccessaryJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.EmptyResult()
	})

	backgroundJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.EmptyResult()
	})

	assert.NoError(t, jm.Schedule(Urgent, urgentJob))
	assert.NoError(t, jm.Schedule(Necessary, neccessaryJob))
	assert.NoError(t, jm.Schedule(Background, backgroundJob))

	count := 0
	jm.Hook(func(ctx context.Context, job priorityqueue.Element[*gojm.Job]) {
		defer func() {
			count++
		}()

		switch count {
		case 0:
			assert.Equal(t, Urgent, job.OriginalPriority())
		case 1:
			assert.Equal(t, Necessary, job.OriginalPriority())
		case 2:
			assert.Equal(t, Background, job.OriginalPriority())
		default:
			assert.FailNow(t, "exceed number of jobs")
		}
	})

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err := jm.RunOne(ctx)
	assert.ErrorIs(t, err, priorityqueue.ErrTimeout)
}

func Test_JobManager_Run1Thread(t *testing.T) {
	jm := gojm.New()
	jm.AddPriority(Urgent)
	jm.AddPriority(Necessary)
	jm.AddPriority(Background)

	urgentJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.EmptyResult()
	})

	neccessaryJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.EmptyResult()
	})

	backgroundJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.EmptyResult()
	})

	assert.NoError(t, jm.Schedule(Urgent, urgentJob))
	assert.NoError(t, jm.Schedule(Necessary, neccessaryJob))
	assert.NoError(t, jm.Schedule(Background, backgroundJob))

	count := 0
	jm.Hook(func(ctx context.Context, job priorityqueue.Element[*gojm.Job]) {
		defer func() {
			count++
		}()

		switch count {
		case 0:
			assert.Equal(t, Urgent, job.OriginalPriority())
		case 1:
			assert.Equal(t, Necessary, job.OriginalPriority())
		case 2:
			assert.Equal(t, Background, job.OriginalPriority())
		default:
			assert.FailNow(t, "exceed number of jobs")
		}
	})

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	err := jm.Run(ctx, 1)
	assert.ErrorIs(t, err, priorityqueue.ErrTimeout)
}

func Test_JobManager_AddJobsAfterRun(t *testing.T) {
	jm := gojm.New()
	jm.AddPriority(Urgent)
	jm.AddPriority(Necessary)
	jm.AddPriority(Background)

	urgentJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		time.Sleep(100 * time.Millisecond)
		return gojm.EmptyResult()
	})

	neccessaryJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		time.Sleep(100 * time.Millisecond)
		return gojm.EmptyResult()
	})

	backgroundJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		time.Sleep(100 * time.Millisecond)
		return gojm.EmptyResult()
	})

	go func() {
		time.Sleep(500 * time.Millisecond)

		assert.NoError(t, jm.Schedule(Background, backgroundJob))
		time.Sleep(20 * time.Millisecond)
		assert.NoError(t, jm.Schedule(Necessary, neccessaryJob))
		time.Sleep(20 * time.Millisecond)
		assert.NoError(t, jm.Schedule(Urgent, urgentJob))
	}()

	count := 0
	jm.Hook(func(ctx context.Context, job priorityqueue.Element[*gojm.Job]) {
		defer func() {
			count++
		}()

		switch count {
		case 0:
			assert.Equal(t, Background, job.OriginalPriority())
		case 1:
			assert.Equal(t, Urgent, job.OriginalPriority())
		case 2:
			assert.Equal(t, Necessary, job.OriginalPriority())
		default:
			assert.FailNow(t, "exceed number of jobs")
		}
	})

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	err := jm.Run(ctx, 1)
	assert.ErrorIs(t, err, priorityqueue.ErrTimeout)
}
