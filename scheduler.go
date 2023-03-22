package scheduler

import (
	"time"

	"github.com/rs/zerolog/log"
)

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
	jobs      jobArray
	tickRate  time.Duration
	isRunning bool
}

// New creates a new scheduler.
func New(tickRate time.Duration) *scheduler {
	log.Info().Int64("tickRate", tickRate.Milliseconds()).Msg("NewScheduler")
	return &scheduler{
		tickRate: tickRate,
		done:     make(chan bool),
	}
}

// Run starts the scheduler.
func (s *scheduler) Run() {
	log.Info().Msg("Scheduler.Run")
	s.isRunning = true
	go s.loop()
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
	var tickNow int64
	for {
		select {
		case <-ticker.C:
			log.Info().Msg("Scheduler.tick")
			timeNow = time.Now().UnixMilli()
			tickDelta = timeNow - timeStart
			tickNow += tickDelta
			s.onTick(tickNow, tickDelta)
		case <-s.done:
			return
		}
	}
}

func (s *scheduler) onTick(t, dt int64) {
	log.Info().Int64("dt", dt).Msg("Scheduler.onTick")
	for _, j := range s.jobs {
		if !j.Valid() {
			log.Info().Msg("Scheduler.onTick: job is not valid. removed")
			s.jobs.Remove(j)
			continue
		}
		if j.Tick(t, dt) {
			log.Info().Msg("Scheduler.onTick: job shall run")
			j.Run()
		}
		if j.Finished() {
			log.Info().Msg("Scheduler.onTick: job finished. removed")
			s.jobs.Remove(j)
		}
	}
}

// CreateJob creates a new job.
func (s *scheduler) CreateJob() *job {
	j := NewJob()
	s.jobs.Append(j)
	return j
}

// Jobs returns the job collection.
func (s *scheduler) Jobs() *jobArray {
	return &s.jobs
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
