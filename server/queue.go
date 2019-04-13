package server

import (
	"context"
	"github.com/kohkimakimoto/hq/hq"
	"sync"
)

type QueueManager struct {
	App         *App
	Queue       chan *hq.Job
	Dispatchers []*Dispatcher
	WorkerWg    *sync.WaitGroup

	// job status
	statusMutex *sync.Mutex
	WaitingJobs map[uint64]*WaitingJob
	RunningJobs map[uint64]*RunningJob
}

func NewQueueManager(app *App) *QueueManager {
	return &QueueManager{
		App:         app,
		Queue:       make(chan *hq.Job, app.Config.Queues),
		Dispatchers: []*Dispatcher{},
		WorkerWg:    &sync.WaitGroup{},
		statusMutex: new(sync.Mutex),
		WaitingJobs: map[uint64]*WaitingJob{},
		RunningJobs: map[uint64]*RunningJob{},
	}
}

func (m *QueueManager) Start() {
	config := m.App.Config

	for i := int64(0); i < config.Dispatchers; i++ {
		d := &Dispatcher{
			manager:    m,
			NumWorkers: 0,
		}
		m.Dispatchers = append(m.Dispatchers, d)

		go d.loop()
	}
}

func (m *QueueManager) Wait() {
	m.WorkerWg.Wait()
}

func (m *QueueManager) EnqueueAsync(job *hq.Job) {
	m.RegisterWaitingJob(job)

	go func() {
		m.Queue <- job
	}()
}

func (m *QueueManager) RegisterRunningJob(job *hq.Job, cancel context.CancelFunc) {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	m.RunningJobs[job.ID] = &RunningJob{
		Job:    job,
		Cancel: cancel,
	}

	// remove waiting jobs
	delete(m.WaitingJobs, job.ID)
}

func (m *QueueManager) RemoveRunningJob(job *hq.Job) {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	delete(m.RunningJobs, job.ID)
}

func (m *QueueManager) RegisterWaitingJob(job *hq.Job) {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	m.WaitingJobs[job.ID] = &WaitingJob{
		Job: job,
	}
}

func (m *QueueManager) UpdateJobStatus(job *hq.Job) *hq.Job {
	if _, ok := m.RunningJobs[job.ID]; ok {
		job.Running = true
	} else if _, ok := m.WaitingJobs[job.ID]; ok {
		job.Waiting = true
	}

	return job
}

func (m *QueueManager) CancelJob(id uint64) error {
	m.statusMutex.Lock()
	defer m.statusMutex.Unlock()

	if rJob, ok := m.RunningJobs[id]; ok {
		rJob.Job.Canceled = true
		rJob.Cancel()
	} else if wJob, ok := m.WaitingJobs[id]; ok {
		wJob.Job.Canceled = true
	}

	return nil
}

type WaitingJob struct {
	Job *hq.Job
}

type RunningJob struct {
	Job    *hq.Job
	Cancel context.CancelFunc
}
