package scheduler

import (
	"time"
)

type TickInfo struct {
	TickCount int64
	TickNow   int64
	TickDelta int64
}

// scheduler is a tick based scheduler.
// A tick is a unit of time that is used to schedule jobs.
// The tick rate is the time between two ticks.
// The tick rate is set when creating a new scheduler.
// The scheduler is started by calling the Run method.
// The scheduler will run until the context is canceled.
// A Job is a function that is executed at a specific time.
// A Job execution will block the scheduler. So make sure that the job is fast.
// The schedule is designed to run in millisecond resolution.
type scheduler struct {
	done      chan bool
	jobs      JobArray
	tickRate  time.Duration
	isRunning bool
	onTick    func(t, dt int64)
	onJobAdd  func(j *Job)
	onJobRm   func(j *Job)
}

// New creates a new scheduler.
func New(tickRate time.Duration) *scheduler {
	return &scheduler{
		tickRate: tickRate,
		done:     make(chan bool),
	}
}

// RunAsync starts the scheduler in a new goroutine.
func (s *scheduler) RunAsync() {
	s.isRunning = true
	go s.loop()
}

// Run starts the scheduler.
func (s *scheduler) Run() {
	s.isRunning = true
	s.loop()
}

// IsRunning returns true if the scheduler is running.
func (s *scheduler) IsRunning() bool {
	return s.isRunning
}

func (s *scheduler) loop() {
	ticker := time.NewTicker(s.tickRate)
	defer func() {
		s.isRunning = false
		ticker.Stop()
	}()
	var timeNow int64
	var tickDelta int64
	timeStart := time.Now().UnixMilli()
	timeLast := timeStart
	var tickNow int64
	for {
		select {
		case <-ticker.C:
			timeNow = time.Now().UnixMilli()

			tickDelta = timeNow - timeLast
			tickNow = timeNow - timeStart
			s.doTick(tickNow, tickDelta)
			timeLast = timeNow
		case <-s.done:
			return
		}
	}
}

// OnTick is called on every tick.
func (s *scheduler) OnTick(f func(t, dt int64)) {
	s.onTick = f
}

func (s *scheduler) doTick(t, dt int64) {
	if s.onTick != nil {
		s.onTick(t, dt)
	}
	for _, j := range s.jobs {
		if !j.Valid() {
			if s.onJobRm != nil {
				s.onJobRm(j)
			}
			s.jobs.Remove(j)
			continue
		}
		if j.Tick(t) {
			j.Run(t)
		}
		if j.Finished() {
			if s.onJobRm != nil {
				s.onJobRm(j)
			}
			s.jobs.Remove(j)
		}
	}
}

// CreateJob creates a new job.
func (s *scheduler) CreateJob() *Job {
	j := NewJob()
	s.jobs.Append(j)
	if s.onJobAdd != nil {
		s.onJobAdd(j)
	}
	return j
}

// Jobs returns the job collection.
func (s *scheduler) Jobs() *JobArray {
	return &s.jobs
}

// OnJobAdd is called when a job is added.
func (s *scheduler) OnJobAdd(f func(j *Job)) {
	s.onJobAdd = f
}

// OnJobRemove is called when a job is removed.
func (s *scheduler) OnJobRemove(f func(j *Job)) {
	s.onJobRm = f
}

// Clear removes all jobs.
func (s *scheduler) Clear() {
	s.Stop()
	s.jobs.Clear()
}

// Stop stops the scheduler.
func (s *scheduler) Stop() {
	s.done <- true
}
