package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	WorkerDefaultUserAgent = fmt.Sprintf("%s/%s", hq.DisplayName, hq.Version)
)

type Dispatcher struct {
	manager    *QueueManager
	NumWorkers int64
}

func (d *Dispatcher) loop() {
	m := d.manager
	logger := m.App.Logger
	config := m.App.Config

	for {
		job := <-m.Queue
		logger.Debugf("dequeue job: %d", job.ID)

		if atomic.LoadInt64(&config.MaxWorkers) <= 0 {
			// sync
			d.dispatch(job)
		} else if atomic.LoadInt64(&d.NumWorkers) < atomic.LoadInt64(&config.MaxWorkers) {
			// async
			d.dispatchAsync(job)
		} else {
			// sync
			d.dispatch(job)
		}
	}
}

func (d *Dispatcher) dispatchAsync(job *hq.Job) {
	manager := d.manager

	manager.WorkerWg.Add(1)
	atomic.AddInt64(&d.NumWorkers, 1)

	go func() {
		defer func() {
			manager.WorkerWg.Done()
			atomic.AddInt64(&d.NumWorkers, -1)
		}()

		d.work(job)
	}()
}

func (d *Dispatcher) dispatch(job *hq.Job) {
	manager := d.manager

	manager.WorkerWg.Add(1)
	atomic.AddInt64(&d.NumWorkers, 1)
	defer func() {
		manager.WorkerWg.Done()
		atomic.AddInt64(&d.NumWorkers, -1)
	}()

	d.work(job)
}

func (d *Dispatcher) work(job *hq.Job) {
	manager := d.manager
	app := manager.App
	logger := app.Logger
	store := app.Store

	logger.Infof("job: %d working", job.ID)

	// worker error
	var err error

	// the terminating logic
	defer func() {
		logger.Infof("job: %d worked", job.ID)
		logger.Debugf("job: %d closing", job.ID)

		// Update result status (success, failure or canceled).
		// If the evaluator has an error, write it to the output buf.
		if err != nil {
			logger.Errorf("worker error: %v", err)
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

	// worker context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// keep running job status.
	manager.RegisterRunningJob(job, cancel)
	defer manager.RemoveRunningJob(job)

	if job.Canceled {
		return
	}

	// Truncate millisecond. It is compatible time for katsubushi ID generator time stamp.
	now := time.Now().UTC().Truncate(time.Millisecond)
	// update startedAt
	job.StartedAt = &now
	if e := store.UpdateJob(job); e != nil {
		logger.Error(e)
	}

	// worker
	err = d.runHttpWorker(job, ctx)
}

func (d *Dispatcher) runHttpWorker(job *hq.Job, ctx context.Context) error {
	// worker
	req, err := http.NewRequest(
		"POST",
		job.URL,
		bytes.NewReader(job.Payload),
	)
	if err != nil {
		return errors.Wrap(err, "failed to create new request")
	}

	// set context
	req = req.WithContext(ctx)

	// common headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("User-Agent", WorkerDefaultUserAgent)
	req.Header.Add("X-Hq-Job-Id", fmt.Sprintf("%d", job.ID))

	// job specific headers
	for k, v := range job.Headers {
		req.Header.Add(k, v)
	}

	// http client
	client := &http.Client{
		Timeout: time.Duration(job.Timeout) * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do http request")
	}
	defer resp.Body.Close()

	statusCode := resp.StatusCode

	job.StatusCode = &statusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read http response body")
	}
	job.Output = string(body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(http.StatusText(resp.StatusCode))
	}

	return nil
}
