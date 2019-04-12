package server

import (
	"context"
	"github.com/kohkimakimoto/hq/structs"
	"sync"
)

type QueueManager struct {
	App         *App
	Queue       chan *structs.Job
	Dispatchers []*Dispatcher
	WorkerWg    *sync.WaitGroup
	mutex       *sync.Mutex
	runningJobs map[uint64]*RunningJob
}

func NewQueueManager(app *App) *QueueManager {
	return &QueueManager{
		App:         app,
		Queue:       make(chan *structs.Job, app.Config.Queues),
		Dispatchers: []*Dispatcher{},
		WorkerWg:    &sync.WaitGroup{},
		mutex:       new(sync.Mutex),
		runningJobs: map[uint64]*RunningJob{},
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

func (m *QueueManager) EnqueueAsync(job *structs.Job) {
	go func() {
		m.Queue <- job
	}()
}

func (m *QueueManager) SetRunningJob(job *structs.Job, cancel context.CancelFunc) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.runningJobs[job.ID] = &RunningJob{
		Job:    job,
		Cancel: cancel,
	}
}

func (m *QueueManager) RemoveRunningJob(job *structs.Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	delete(m.runningJobs, job.ID)
}

func (m *QueueManager) UpdateRunningStatus(job *structs.Job) *structs.Job {
	if _, ok := m.runningJobs[job.ID]; ok {
		job.Running = true
	} else {
		job.Running = false
	}
	return job
}

type RunningJob struct {
	Job    *structs.Job
	Cancel context.CancelFunc
}
