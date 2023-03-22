# Job scheduler

A job scheduler is a system that is used to schedule jobs to run at an interval. It is used to schedule jobs that are not time critical and can be run at a later time. 

The scheduler is not ware of any time or time location. It is based on a time interval. The time interval is the time between the current time and the time the job is scheduled to run. All time presented here are in milliseconds.

When a job is finished the job will automatically be removed from the scheduler. If you want to run the job again you need to create a new job.

The scheduler reschedules the job to run at the next interval, if it is repeated or runs forever. If a job is scheduled in the past it will run immediately.

## Installation

```bash
go get github.com/apigear-io/scheduler
```

## Usage

```go
s := scheduler.NewScheduler(time.Millisecond * 10)
s.CreateJob().Every(time.Millisecond * 100).Do(func() {
    fmt.Println("Hello World")
})
s.Start()
```

