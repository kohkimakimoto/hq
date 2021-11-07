package server

import (
	"context"
	"sync"

	"github.com/kohkimakimoto/hq/internal/structs"
)

// QueueManager provides a queue data structure and manages the queued job status.
type QueueManager struct {
	Queue chan *structs.Job

	// properties for job status
	mutex       sync.RWMutex
	waitingJobs map[uint64]*WaitingJob
	runningJobs map[uint64]*RunningJob
}

func NewQueueManager(queueSize int64) *QueueManager {
	return &QueueManager{
		Queue:       make(chan *structs.Job, queueSize),
		waitingJobs: map[uint64]*WaitingJob{},
		runningJobs: map[uint64]*RunningJob{},
	}
}

func (m *QueueManager) EnqueueAsync(job *structs.Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// set the job as a waiting job
	m.waitingJobs[job.ID] = &WaitingJob{
		Job: job,
	}

	// As enqueuing jobs asynchronously, it does NOT guarantee the order of the jobs.
	go func() {
		m.Queue <- job
	}()
}

func (m *QueueManager) Dequeue() *structs.Job {
	return <-m.Queue
}

func (m *QueueManager) RegisterRunningJob(job *structs.Job, cancel context.CancelFunc) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.runningJobs[job.ID] = &RunningJob{
		Job:    job,
		Cancel: cancel,
	}

	// remove waiting jobs
	delete(m.waitingJobs, job.ID)
}

func (m *QueueManager) RemoveRunningJob(job *structs.Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.runningJobs, job.ID)
}

func (m *QueueManager) CancelJob(id uint64) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if rJob, ok := m.runningJobs[id]; ok {
		rJob.Job.Canceled = true
		rJob.Cancel()
	} else if wJob, ok := m.waitingJobs[id]; ok {
		wJob.Job.Canceled = true
	}
}

func (m *QueueManager) LoadJobStatus(job *structs.Job) *structs.Job {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	if rJob, ok := m.runningJobs[job.ID]; ok {
		job.Running = true
		job.Canceled = rJob.Job.Canceled
	} else if wJob, ok := m.waitingJobs[job.ID]; ok {
		job.Waiting = true
		job.Canceled = wJob.Job.Canceled
	}

	return job
}

func (m *QueueManager) NumJobsWaiting() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.waitingJobs)
}

func (m *QueueManager) NumJobsRunning() int {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	return len(m.runningJobs)
}

func (m *QueueManager) NumJobsInQueue() int {
	return len(m.Queue)
}

type WaitingJob struct {
	Job *structs.Job
}

type RunningJob struct {
	Job    *structs.Job
	Cancel context.CancelFunc
}
