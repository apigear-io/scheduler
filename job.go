package scheduler

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type job struct {
	// id is the unique identifier of the job.
	id string
	// tickNow is the current tick of the scheduler.
	tickNow int64
	// tickRun is the tick when the job should run next.
	tickRun int64
	// interval is the interval in ticks between two runs.
	interval int64
	// repeat is the number of times the job should be repeated.
	repeat int64
	// forever is true if the job should be repeated forever.
	forever bool
	// tags is a map of tags of the job.
	tags map[string]bool
	fn   func()
}

func NewJob() *job {
	return &job{
		id:       uuid.NewString(),
		tickRun:  0,
		interval: 0,
		repeat:   0,
		forever:  false,
		tags:     make(map[string]bool),
		fn:       nil,
	}
}

func (j *job) Tick(t, dt int64) bool {
	j.tickNow = t
	return t >= j.tickRun
}

func (j *job) Run() {
	if !j.forever && j.repeat == 0 {
		return
	}
	if j.repeat > 0 {
		j.repeat--
	}
	j.fn()
	j.tickRun += j.interval
}

// Finished returns true if the job has finished running.
func (j *job) Finished() bool {
	return !j.forever && j.repeat == 0
}

// Valid returns true if the job is valid.
func (j *job) Valid() bool {
	return j.fn != nil && j.tickRun > 0
}

// At runs the job at the given tick.
// The job will be run at the next tick if the given tick is in the past.
func (j *job) At(t int64) *job {
	j.tickRun = t
	if j.repeat == 0 && !j.forever {
		j.repeat = 1
		j.forever = false
	}
	return j
}

// Every runs the job every interval.
func (j *job) Every(interval time.Duration) *job {
	j.interval = interval.Milliseconds()
	j.tickRun = j.tickNow + j.interval
	if j.repeat == 0 {
		j.forever = true
	}
	return j
}

// After runs the job after interval.
func (j *job) After(interval time.Duration) *job {
	j.interval = interval.Milliseconds()
	j.tickRun = j.tickNow + j.interval
	if j.repeat == 0 && !j.forever {
		j.repeat = 1
		j.forever = false
	}
	return j
}

// Repeat sets the number of times the job should be repeated.
func (j *job) Repeat(n int64) *job {
	j.repeat = n
	j.forever = false
	return j
}

// Forever sets the job to run forever.
func (j *job) Forever() *job {
	j.repeat = 0
	j.forever = true
	return j
}

// Do sets the function of the job.
func (j *job) Do(fn func()) *job {
	j.fn = fn
	return j
}

// Now runs the job now.
func (j *job) Now() {
	j.tickRun = j.tickNow
	j.Run()
}

// HasTag returns true if the job has the given tag.
func (j *job) HasTag(tag string) bool {
	return j.tags[tag]
}

// Tags returns a slice of all tags of the job.
func (j *job) TagList() []string {
	tags := make([]string, 0, len(j.tags))
	for tag := range j.tags {
		tags = append(tags, tag)
	}
	return tags
}

// Tag returns a new job with the given tag.
func (j *job) Tag(tag ...string) *job {
	for _, tag := range tag {
		j.tags[tag] = true
	}
	return j
}

// Untag returns a new job without the given tag.
func (j *job) Untag(tag ...string) *job {
	for _, tag := range tag {
		delete(j.tags, tag)
	}
	return j
}

// String returns a string representation of the job.
func (j *job) String() string {
	return fmt.Sprintf("Job{id:%s, tickNow:%d, tickRun:%d, interval:%d, repeat:%d, forever:%t, tags:%v}",
		j.id, j.tickNow, j.tickRun, j.interval, j.repeat, j.forever, j.tags)
}
