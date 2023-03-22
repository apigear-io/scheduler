package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionCreate(t *testing.T) {
	c := jobArray{}
	j := c.CreateJob()
	assert.Equal(t, 1, len(c))
	assert.Equal(t, j, c[0])
}

func TestCollectionRemove(t *testing.T) {
	c := jobArray{}
	j := c.CreateJob()
	c.Remove(j)
	assert.Equal(t, 0, len(c))
}

func TestCollectionClear(t *testing.T) {
	c := jobArray{}
	c.CreateJob()
	c.CreateJob()
	c.CreateJob()
	c.Clear()
	assert.Equal(t, 0, len(c))
}
