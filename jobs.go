package scheduler

// JobArray is aan array of jobs.
type JobArray []*Job

// Len returns the number of jobs in the array.
func (c *JobArray) Len() int {
	return len(*c)
}

// Append adds the job to the array.
func (c *JobArray) Append(j *Job) {
	*c = append(*c, j)
}

// Remove removes the job from the array.
func (c *JobArray) Remove(j *Job) {
	for i, e := range *c {
		if e == j {
			*c = append((*c)[:i], (*c)[i+1:]...)
			return
		}
	}
}

// CreateJob creates a new job and adds it to the array.
func (c *JobArray) CreateJob() *Job {
	j := NewJob()
	c.Append(j)
	return j
}

// Clear removes all jobs from the array.
func (c *JobArray) Clear() {
	*c = (*c)[:0]
}
