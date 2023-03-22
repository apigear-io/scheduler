package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSchedulerCreateJob(t *testing.T) {
	s := New(time.Millisecond * 10)
	s.CreateJob()
	assert.Equal(t, 1, s.Jobs().Len())
}

func TestSchedulerRun(t *testing.T) {
	s := New(time.Millisecond * 10)
	s.CreateJob()
	s.RunAsync()
	assert.Equal(t, true, s.IsRunning())
	s.Stop()
	assert.Equal(t, false, s.IsRunning())
}

func TestSchedulerRunAfter(t *testing.T) {
	s := New(time.Millisecond * 10)
	hasRun := false
	s.CreateJob().Do(func(t, dt int64) {
		hasRun = true
	}).After(time.Millisecond * 10)
	s.RunAsync()
	assert.Equal(t, true, s.IsRunning())
	time.Sleep(time.Millisecond * 20)
	assert.Equal(t, true, hasRun)
	s.Stop()
	assert.Eventually(t, func() bool {
		return !s.IsRunning()
	}, time.Millisecond*100, time.Millisecond*10)
}

func TestSchedulerRunEvery(t *testing.T) {
	s := New(time.Millisecond * 10)
	runCount := 0
	s.CreateJob().Do(func(t, dt int64) {
		runCount++
	}).Every(time.Millisecond * 10)
	s.RunAsync()
	assert.Equal(t, true, s.IsRunning())
	time.Sleep(time.Millisecond * 25)
	assert.Equal(t, 2, runCount)
	s.Stop()
	assert.Eventually(t, func() bool {
		return !s.IsRunning()
	}, time.Millisecond*100, time.Millisecond*10)
}

func TestSchedulerRunRepeat(t *testing.T) {
	s := New(time.Millisecond * 10)
	runCount := 0
	s.CreateJob().Do(func(t, dt int64) {
		runCount++
	}).Every(time.Millisecond * 10).Repeat(2)
	s.RunAsync()
	assert.Equal(t, true, s.IsRunning())
	time.Sleep(time.Millisecond * 25)
	assert.Equal(t, 2, runCount)
	time.Sleep(time.Millisecond * 25)
	assert.Equal(t, 2, runCount)
	s.Stop()
	assert.Eventually(t, func() bool {
		return !s.IsRunning()
	}, time.Millisecond*100, time.Millisecond*10)
}

func TestSchedulerRunForever(t *testing.T) {
	s := New(time.Millisecond * 10)
	runCount := 0
	s.CreateJob().Do(func(t, dt int64) {
		runCount++
	}).Every(time.Millisecond * 1).Forever()
	s.RunAsync()
	assert.Equal(t, true, s.IsRunning())
	time.Sleep(time.Millisecond * 25)
	assert.Equal(t, 2, runCount)
	time.Sleep(time.Millisecond * 30)
	assert.Equal(t, 5, runCount)
	s.Stop()
}
