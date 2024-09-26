package gojm_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xybor-x/gojm"
)

func Test_Job_Execute(t *testing.T) {
	job := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.Result(1)
	})

	ctx := context.Background()
	result := job.Execute(ctx)
	assert.Equal(t, 1, result.GetInt(nil))
}

func Test_Job_GetResult(t *testing.T) {
	job := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.Result(1)
	})

	ctx := context.Background()
	job.Execute(ctx)

	result := job.GetResult()
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.GetInt(nil))
}

func Test_Job_WaitResult(t *testing.T) {
	job := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.Result(1)
	})

	ctx := context.Background()
	go func(ctx context.Context) {
		time.Sleep(500 * time.Millisecond)
		job.Execute(ctx)
	}(ctx)

	// Run first time.
	ctx = context.Background()
	result := job.WaitResult(ctx)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.GetInt(nil))

	// Run second time, expect that it is the same as previous time.
	result = job.WaitResult(ctx)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.GetInt(nil))
}

func Test_Job_WaitResultTimeout(t *testing.T) {
	job := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		return gojm.Result(1)
	})

	ctx := context.Background()
	go func(ctx context.Context) {
		time.Sleep(500 * time.Millisecond)
		job.Execute(ctx)
	}(ctx)

	ctx = context.Background()
	ctx, cancel := context.WithTimeout(ctx, 100*time.Millisecond)
	defer cancel()

	result := job.WaitResult(ctx)
	assert.Nil(t, result)
}
