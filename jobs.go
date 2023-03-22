package scheduler

// jobArray is aan array of jobs.
type jobArray []*job

// Len returns the number of jobs in the array.
func (c *jobArray) Len() int {
	return len(*c)
}

// Append adds the job to the array.
func (c *jobArray) Append(j *job) {
	*c = append(*c, j)
}

// Remove removes the job from the array.
func (c *jobArray) Remove(j *job) {
	for i, e := range *c {
		if e == j {
			*c = append((*c)[:i], (*c)[i+1:]...)
			return
		}
	}
}

// CreateJob creates a new job and adds it to the array.
func (c *jobArray) CreateJob() *job {
	j := NewJob()
	c.Append(j)
	return j
}

// Clear removes all jobs from the array.
func (c *jobArray) Clear() {
	*c = (*c)[:0]
}
