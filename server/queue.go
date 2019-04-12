package server

import (
	"context"
	"github.com/kohkimakimoto/hq/structs"
	"sync"
)

type QueueManager struct {
	app         *App
	queue       chan *structs.Job
	dispatchers []*Dispatcher
	mutex       *sync.Mutex
	runningJobs map[uint64]*RunningJob
	wg          *sync.WaitGroup
}

func NewQueueManager(app *App) *QueueManager {
	return &QueueManager{
		app:         app,
		queue:       make(chan *structs.Job, app.Config.Queues),
		dispatchers: []*Dispatcher{},
		mutex:       new(sync.Mutex),
		runningJobs: map[uint64]*RunningJob{},
		wg:          &sync.WaitGroup{},
	}
}

func (m *QueueManager) Start() {
	config := m.app.Config

	for i := int64(0); i < config.Dispatchers; i++ {
		d := &Dispatcher{
			manager:    m,
			numWorkers: 0,
		}
		m.dispatchers = append(m.dispatchers, d)

		go d.loop()
	}
}

func (m *QueueManager) Wait() {
	m.wg.Wait()
}

func (m *QueueManager) Enqueue(job *structs.Job) {
	m.queue <- job
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

func (m *QueueManager) SetRunningStatus(job *structs.Job) *structs.Job {
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
