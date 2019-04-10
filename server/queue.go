package server

import (
	"bytes"
	"fmt"
	"github.com/kohkimakimoto/hq/structs"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

type QueueManager struct {
	app           *App
	queue         chan *structs.Job
	dispatchers   []*Dispatcher
	numWorkersAll int64
	mutex       *sync.Mutex
	runningJobs    map[uint64]*structs.Job
}

func NewQueueManager(app *App) *QueueManager {
	return &QueueManager{
		app:           app,
		queue:         make(chan *structs.Job, app.Config.Queues),
		dispatchers:   []*Dispatcher{},
		numWorkersAll: 0,
		mutex:       new(sync.Mutex),
		runningJobs: map[uint64]*structs.Job{},
	}
}

func (m *QueueManager) Start() {
	config := m.app.Config

	for i := int64(0); i < config.Dispatchers; i++ {
		d := &Dispatcher{
			manager:    m,
			numWorkers: 0,
		}

		go d.loop()
	}
}

func (m *QueueManager) Enqueue(job *structs.Job) {
	m.queue <- job
}

func (m *QueueManager) SetRunningJob(job *structs.Job) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.runningJobs[job.ID] = job
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

type Dispatcher struct {
	manager    *QueueManager
	numWorkers int64
}

func (d *Dispatcher) loop() {
	m := d.manager
	logger := m.app.Logger
	config := m.app.Config

	for {
		job := <-m.queue
		logger.Debugf("dequeue job: %d", job.ID)

		if atomic.LoadInt64(&config.MaxWorkers) <= 0 {
			// sync
			d.work(job)
		} else if atomic.LoadInt64(&d.numWorkers) < atomic.LoadInt64(&config.MaxWorkers) {
			// async
			atomic.AddInt64(&d.numWorkers, 1)
			atomic.AddInt64(&m.numWorkersAll, 1)

			go func(job *structs.Job) {
				d.work(job)
				atomic.AddInt64(&d.numWorkers, -1)
				atomic.AddInt64(&m.numWorkersAll, -1)
			}(job)
		} else {
			// sync
			d.work(job)
		}
	}
}

func (d *Dispatcher) work(job *structs.Job) {
	manager := d.manager
	app := d.manager.app
	logger := app.Logger
	store := app.Store

	logger.Infof("job: %d working", job.ID)

	// worker error
	var err error

	// the terminating logic
	defer func() {
		logger.Infof("job: %d worked", job.ID)
		logger.Debugf("job: %d closing", job.ID)

		// Update result status (success or failure).
		// If the evaluator has an error, write it to the output buf.
		if err != nil {
			job.Success = false
			job.Failure = true
			job.Err = err.Error()
		} else {
			job.Success = true
			job.Failure = false
		}

		// Truncate millisecond. It is compatible time for katsubushi ID generator time stamp.
		now := time.Now().UTC().Truncate(time.Millisecond)

		// update finishedAt
		job.FinishedAt = &now

		if e := store.UpdateJob(job); e != nil {
			logger.Error(e)
		}

		logger.Debugf("job: %d closed", job.ID)
	}()

	// change status running.
	manager.SetRunningJob(job)
	defer manager.RemoveRunningJob(job)

	// worker
	client := &http.Client{
		Timeout: time.Duration(job.Timeout) * time.Second,
	}
	req, err := http.NewRequest(
		"POST",
		job.URL,
		bytes.NewReader(job.Payload),
	)
	if err != nil {
		return
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	job.StatusCode = resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	job.Output = string(body)

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf(http.StatusText(resp.StatusCode))
		return
	}
}
