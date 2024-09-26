package gojm_test

import (
	"context"
	"fmt"
	"time"

	"github.com/xybor-x/gojm"
)

func ExampleJobManager() {
	// Setup priority
	jm := gojm.New()
	jm.AddPriority(Urgent)
	jm.AddPriority(Necessary)
	jm.AddPriority(Background)

	backgroundJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
		fmt.Println("background done")
		return gojm.Result("background")
	})

	// In other goroutines, schedule your jobs.
	go func() {
		time.Sleep(20 * time.Millisecond)

		jm.Schedule(Background, backgroundJob)

		time.Sleep(20 * time.Millisecond)

		jm.Schedule(Urgent, gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
			fmt.Println("urgent done")
			return gojm.EmptyResult()
		}))

		time.Sleep(20 * time.Millisecond)

		jm.Schedule(Necessary, gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
			fmt.Println("necessary done")
			return gojm.EmptyResult()
		}))
	}()

	// Because this is an example, we need to stop this function after one
	// second.
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 1000*time.Millisecond)
	defer cancel()

	// Run the manager with the ability of handling 2 jobs in a time. In
	// reality, you could run this function forever.
	jm.Run(ctx, 2)

	// You also wait for the result of a job.
	result := backgroundJob.WaitResult(ctx)
	fmt.Println("Result of background job:", result.Get(nil))

	// Output:
	// background done
	// urgent done
	// necessary done
	// Result of background job: background
}
