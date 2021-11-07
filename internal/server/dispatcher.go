package server

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"

	"github.com/kohkimakimoto/hq/internal/structs"
	"github.com/kohkimakimoto/hq/internal/version"
)

var (
	WorkerDefaultUserAgent = fmt.Sprintf("HQ/%s", version.Version)
)

// Dispatcher contains multiple workers and dispatches jobs from the queue to the workers.
type Dispatcher struct {
	queueManager      *QueueManager
	store             *Store
	logger            echo.Logger
	httpClientFactory func() *http.Client
	workerWg          sync.WaitGroup
	maxWorkers        int64
	numWorkers        int64
}

func defaultHttpClientFactory() *http.Client {
	return &http.Client{}
}

func (d *Dispatcher) EventLoop() {
	for {
		job := d.queueManager.Dequeue()
		d.logger.Debugf("dequeue job: %d", job.ID)

		if atomic.LoadInt64(&d.maxWorkers) <= 0 {
			// sync
			d.dispatch(job)
		} else if atomic.LoadInt64(&d.numWorkers) < atomic.LoadInt64(&d.maxWorkers) {
			// async
			d.dispatchAsync(job)
		} else {
			// sync
			d.dispatch(job)
		}
	}
}

func (d *Dispatcher) dispatchAsync(job *structs.Job) {
	d.workerWg.Add(1)
	atomic.AddInt64(&d.numWorkers, 1)
	go func() {
		defer func() {
			d.workerWg.Done()
			atomic.AddInt64(&d.numWorkers, -1)
		}()
		d.work(job)
	}()
}

func (d *Dispatcher) dispatch(job *structs.Job) {
	d.workerWg.Add(1)
	atomic.AddInt64(&d.numWorkers, 1)
	defer func() {
		d.workerWg.Done()
		atomic.AddInt64(&d.numWorkers, -1)
	}()
	d.work(job)
}

func (d *Dispatcher) work(job *structs.Job) {
	d.logger.Infof("job: %d started working", job.ID)

	// worker error
	var err error

	// the terminating logic
	defer func() {
		d.logger.Infof("job: %d finished working", job.ID)
		d.logger.Debugf("job: %d closing", job.ID)

		// Update result status (success, failure or canceled).
		// If the evaluator has an error, write it to the output buf.
		if err != nil {
			d.logger.Errorf("worker error: %v", err)
			job.Success = false
			job.Failure = true
			job.Err = err.Error()
		} else {
			job.Success = true
			job.Failure = false
		}

		// Truncate millisecond. It is compatible time for katsubushi ID generator timestamp.
		now := time.Now().UTC().Truncate(time.Millisecond)
		// update finishedAt
		job.FinishedAt = &now
		if e := d.store.UpdateJob(job); e != nil {
			d.logger.Error(e)
		}

		d.logger.Debugf("job: %d closed", job.ID)
	}()

	// worker context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// make job status running.
	d.queueManager.RegisterRunningJob(job, cancel)
	defer d.queueManager.RemoveRunningJob(job)

	if job.Canceled {
		return
	}

	// Truncate millisecond. It is compatible time for katsubushi ID generator timestamp.
	now := time.Now().UTC().Truncate(time.Millisecond)
	// update startedAt
	job.StartedAt = &now
	if e := d.store.UpdateJob(job); e != nil {
		d.logger.Error(e)
	}

	// run worker
	err = d.runHttpWorker(ctx, job)
}

func (d *Dispatcher) runHttpWorker(ctx context.Context, job *structs.Job) error {
	var reqBody io.Reader
	if job.Payload != nil && !bytes.Equal(job.Payload, []byte("null")) {
		reqBody = bytes.NewReader(job.Payload)
	}

	// worker
	req, err := http.NewRequest(
		"POST",
		job.URL,
		reqBody,
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
	client := d.httpClientFactory()
	client.Timeout = time.Duration(job.Timeout) * time.Second

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

// NumWorkers returns the number of workers that are working now.
func (d *Dispatcher) NumWorkers() int64 {
	return atomic.LoadInt64(&d.numWorkers)
}

func (d *Dispatcher) Wait() {
	for {
		// wait for all queued jobs to finish
		numWaitingJobs := d.queueManager.NumJobsWaiting()
		numRunningJobs := d.queueManager.NumJobsRunning()
		numQueuedJobs := d.queueManager.NumJobsInQueue()
		numWorkers := d.NumWorkers()

		if numWaitingJobs == 0 && numRunningJobs == 0 && numQueuedJobs == 0 && numWorkers == 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	// wait for all workers to finish
	d.workerWg.Wait()
}
