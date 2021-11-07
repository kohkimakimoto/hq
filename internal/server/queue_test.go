package server

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kohkimakimoto/hq/internal/structs"
)

func TestQueueManager_EnqueueAsync(t *testing.T) {
	m := NewQueueManager(10)
	for i := uint64(0); i < 20; i++ {
		m.EnqueueAsync(&structs.Job{
			ID:   i,
			Name: fmt.Sprintf("job-%d", i),
		})
	}

	// wait for jobs to be enqueued fully
	for m.NumJobsInQueue() < 10 {
		time.Sleep(time.Millisecond * 1)
	}

	// check queue status
	assert.Equal(t, 10, m.NumJobsInQueue())
	assert.Equal(t, 20, m.NumJobsWaiting())
	assert.Equal(t, 0, m.NumJobsRunning())
}

func TestQueueManager_Dequeue(t *testing.T) {
	m := NewQueueManager(10)
	for i := uint64(0); i < 10; i++ {
		m.EnqueueAsync(&structs.Job{
			ID:   i,
			Name: fmt.Sprintf("job-%d", i),
		})
	}

	// wait for jobs to be enqueued fully
	for m.NumJobsInQueue() < 10 {
		time.Sleep(time.Millisecond * 1)
	}

	// check queue status
	assert.Equal(t, 10, m.NumJobsInQueue())
	assert.Equal(t, 10, m.NumJobsWaiting())
	assert.Equal(t, 0, m.NumJobsRunning())

	// dequeue the first job
	_ = m.Dequeue()

	// check queue status
	assert.Equal(t, 9, m.NumJobsInQueue())
	assert.Equal(t, 10, m.NumJobsWaiting()) // still waiting status
	assert.Equal(t, 0, m.NumJobsRunning())
}

func TestQueueManager_RegisterRunningJob(t *testing.T) {
	m := NewQueueManager(10)
	for i := uint64(0); i < 10; i++ {
		m.EnqueueAsync(&structs.Job{
			ID:   i,
			Name: fmt.Sprintf("job-%d", i),
		})
	}

	// wait for jobs to be enqueued fully
	for m.NumJobsInQueue() < 10 {
		time.Sleep(time.Millisecond * 1)
	}

	// check queue status
	assert.Equal(t, 10, m.NumJobsInQueue())
	assert.Equal(t, 10, m.NumJobsWaiting())
	assert.Equal(t, 0, m.NumJobsRunning())

	// dequeue the first job
	job := m.Dequeue()
	m.RegisterRunningJob(job, func() {})

	// check queue status
	assert.Equal(t, 9, m.NumJobsInQueue())
	assert.Equal(t, 9, m.NumJobsWaiting())
	assert.Equal(t, 1, m.NumJobsRunning())
}

func TestQueueManager_RemoveRunningJob(t *testing.T) {
	m := NewQueueManager(10)
	for i := uint64(0); i < 10; i++ {
		m.EnqueueAsync(&structs.Job{
			ID:   i,
			Name: fmt.Sprintf("job-%d", i),
		})
	}

	// wait for jobs to be enqueued fully
	for m.NumJobsInQueue() < 10 {
		time.Sleep(time.Millisecond * 1)
	}

	// check queue status
	assert.Equal(t, 10, m.NumJobsInQueue())
	assert.Equal(t, 10, m.NumJobsWaiting())
	assert.Equal(t, 0, m.NumJobsRunning())

	// dequeue the first job
	job := m.Dequeue()
	m.RegisterRunningJob(job, func() {})

	// check queue status
	assert.Equal(t, 9, m.NumJobsInQueue())
	assert.Equal(t, 9, m.NumJobsWaiting())
	assert.Equal(t, 1, m.NumJobsRunning())

	// remove the job
	m.RemoveRunningJob(job)

	// check queue status
	assert.Equal(t, 9, m.NumJobsInQueue())
	assert.Equal(t, 9, m.NumJobsWaiting())
	assert.Equal(t, 0, m.NumJobsRunning())
}
