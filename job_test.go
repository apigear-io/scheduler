package scheduler

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestJobTimes(t *testing.T) {
	j := NewJob().Do(func(int64, int64) {}).Repeat(3)
	assert.Equal(t, int64(3), j.repeat)
	assert.Equal(t, false, j.forever)
	assert.Equal(t, int64(0), j.interval)
	assert.Equal(t, int64(0), j.tickRun)
	assert.Equal(t, false, j.Valid())
}

func TestJobInterval(t *testing.T) {
	j := NewJob().Do(func(int64, int64) {}).After(time.Second)
	assert.Equal(t, int64(1000), j.interval)
	assert.Equal(t, false, j.forever)
	assert.Equal(t, int64(1), j.repeat)
	assert.Equal(t, int64(1000), j.tickRun)
	assert.Equal(t, true, j.Valid())
}

func TestJobForever(t *testing.T) {
	j := NewJob().Do(func(int64, int64) {}).Every(time.Second)
	assert.Equal(t, int64(1000), j.interval)
	assert.Equal(t, true, j.forever)
	assert.Equal(t, int64(0), j.repeat)
	assert.Equal(t, int64(1000), j.tickRun)
	assert.Equal(t, true, j.Valid())
}

func TestJobTag(t *testing.T) {
	j := NewJob().Do(func(int64, int64) {}).Tag("tag1", "tag2")
	assert.Equal(t, true, j.tags["tag1"])
	assert.Equal(t, true, j.tags["tag2"])
}

func TestJobUntag(t *testing.T) {
	j := NewJob().Do(func(int64, int64) {}).Tag("tag1", "tag2").Untag("tag1")
	assert.Equal(t, false, j.tags["tag1"])
	assert.Equal(t, true, j.tags["tag2"])
}
func TestJobNotValid(t *testing.T) {
	j := NewJob()
	assert.Equal(t, false, j.Valid())
	j.Do(nil)
	assert.Equal(t, false, j.Valid())
	j.Do(func(int64, int64) {})
	assert.Equal(t, false, j.Valid())
	j.After(time.Second)
	assert.Equal(t, true, j.Valid())
}
