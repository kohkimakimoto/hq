package server

import (
	"context"
	"github.com/kohkimakimoto/hq/hq"
	"sync"
)

type QueueManager struct {
	App              *App
	Queue            chan *hq.Job
	Dispatchers      []*Dispatcher
	WorkerWg         *sync.WaitGroup
	WaitingJobs      map[uint64]*WaitingJob
	waitingJobsMutex *sync.Mutex
	RunningJobs      map[uint64]*RunningJob
	runningJobsMutex *sync.Mutex
}

func NewQueueManager(app *App) *QueueManager {
	return &QueueManager{
		App:              app,
		Queue:            make(chan *hq.Job, app.Config.Queues),
		Dispatchers:      []*Dispatcher{},
		WorkerWg:         &sync.WaitGroup{},
		WaitingJobs:      map[uint64]*WaitingJob{},
		waitingJobsMutex: new(sync.Mutex),
		RunningJobs:      map[uint64]*RunningJob{},
		runningJobsMutex: new(sync.Mutex),
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
	m.runningJobsMutex.Lock()
	defer m.runningJobsMutex.Unlock()

	m.RunningJobs[job.ID] = &RunningJob{
		Job:    job,
		Cancel: cancel,
	}

	m.RemoveWaitingJob(job)
}

func (m *QueueManager) RemoveRunningJob(job *hq.Job) {
	m.runningJobsMutex.Lock()
	defer m.runningJobsMutex.Unlock()

	delete(m.RunningJobs, job.ID)
}

func (m *QueueManager) RegisterWaitingJob(job *hq.Job) {
	m.waitingJobsMutex.Lock()
	defer m.waitingJobsMutex.Unlock()

	m.WaitingJobs[job.ID] = &WaitingJob{
		Job:      job,
		Canceled: false,
	}
}

func (m *QueueManager) RemoveWaitingJob(job *hq.Job) {
	m.waitingJobsMutex.Lock()
	defer m.waitingJobsMutex.Unlock()

	delete(m.WaitingJobs, job.ID)
}

func (m *QueueManager) UpdateJobStatus(job *hq.Job) *hq.Job {
	if _, ok := m.RunningJobs[job.ID]; ok {
		job.Running = true
	} else if _, ok := m.WaitingJobs[job.ID]; ok {
		job.Waiting = true
	}

	return job
}

type WaitingJob struct {
	Job      *hq.Job
	Canceled bool
}

type RunningJob struct {
	Job    *hq.Job
	Cancel context.CancelFunc
}
