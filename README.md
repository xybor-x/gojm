# Introduction

A thread-safe and reliable priority-based job manager.

# Job

In `gojm`, `job` is only a wrapper of function. You can define a job as below:
```golang
job := gojm.NewJob(function(ctx context.Context) *gojm.JobResult {
    fmt.Println("The job started")

    time.Sleep(time.Second)

    fmt.Println("The job completed")

    return nil
})
```

You can execute the job like you call the function. Note that the method
`Exec()` should be called only once. The program will panics if you call it
again.

```golang
ctx := context.Background()
result := job.Exec(ctx)
```

You also get the result with blocking mode and non-blocking mode after the job
already completed. If the job has not completed, the result is nil.

```golang
// Non-blocking mode
result := job.GetResult()
```

```golang
// Blocking mode
ctx := context.Background()
ctx, cancel := context.WithTimeout(ctx, time.Second) // set timeout as 1 second.
defer cancel()

result := job.WaitResult(ctx)
```

# Job result

If you want to return some values for the job, you can modify JobResult.

```golang
job := gojm.NewJob(function(ctx context.Context) *gojm.JobResult {
    fmt.Println("The job started")

    time.Sleep(time.Second)

    fmt.Println("The job completed")

    return gojm.Result(100)
})

result := job.Exec(ctx)
fmt.Println(result.Get(nil))
// Output:
// 100
```

You also can put many values into JobResult.

```golang
job := gojm.NewJob(function(ctx context.Context) *gojm.JobResult {
    fmt.Println("The job started")

    time.Sleep(time.Second)

    fmt.Println("The job completed")

    result := gojm.EmptyResult()
    result.Set("x", 100)
    result.Set("y", "abc")
})

result := job.Exec(ctx)
fmt.Println(result.Get("x"), result.Get("y"))
// Output:
// 100 abc
```

If you want to return an error, please use `gojm.Err` function.

```golang
job := gojm.NewJob(function(ctx context.Context) *gojm.JobResult {
    fmt.Println("The job started")

    time.Sleep(time.Second)

    fmt.Println("The job completed")

    return gojm.Err(errors.New("something's wrong"))
})

result := job.Exec(ctx)
fmt.Println(result.Err)
// Output:
// something's wrong
```

# Job manager

You can put jobs into a job manager with a priority to execute it in when
possible.

Firstly, you need to create some `Priority`. Every `Priority` has its own value,
the lower value, the higher priority.

```golang
// Urgent is for jobs which need to execute as soon as possible.
var Urgent = gojm.NewPriority("Urgent", 0)

// Necessary is for jobs which can be executed later but also need to be
// completed soon. We set the aging by one minute (after one minute, this job
// will be moved to the higher priority).
var Necessary = gojm.NewPriority("Necessary", 10).WithAging(time.Minute)

// Background is for jobs which can be executed at any time. We must specify we
// don't need an aging for this priority.
var Background = gojm.NewPriority("Background", 1000).WithNoAging()
```

Start a job manager
```golang
jm := gojm.New()

ctx := context.Background()
if err := jm.Run(ctx); err != nil {
    panic(err)
}
```

You can put job into job manager in other goroutines.
```golang
jm.Schedule(Background, gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
    fmt.Println("Job do something")
    return nil
}))
```

Or you can put the job into manager and wait for its result.

```golang
urgentJob := gojm.NewJob(func(ctx context.Context) *gojm.JobResult {
    fmt.Println("Do urgent work")
    return gojm.Result(0)
})

jm.Schedule(Background, urgentJob)

result := urgentJob.WaitResult(ctx)
fmt.Println(result.Get(nil))
// Output:
// 0
```

# Hook

Instead of waiting for the result of each job, you can put a hook function to
handle the result of all completed jobs.

```golang
jm.Hook(func (ctx context.Context, job gojm.JobWrapper) {
    if job.Unwrap().Err != nil {
        log.Printf("level=error priority=%s err=%v", job.OriginalPriority, job.Unwrap().Err)
    } else if result := job.Unwrap().GetResult().Get(nil); result != nil {
        log.Printf("level=info priority=%s result=%v", job.OriginalPriority, result)
    }
})
```
