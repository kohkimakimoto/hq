package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_Status(t *testing.T) {
	job := &Job{}
	job.Waiting = true
	assert.Equal(t, JobStatusWaiting, job.Status())
}
