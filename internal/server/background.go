package server

import (
	"sync"
	"time"

	"github.com/kayac/go-katsubushi"
	"github.com/labstack/echo/v4"
)

type BackgroundCleaner struct {
	logger       echo.Logger
	queueManager *QueueManager
	store        *Store
	jobLifetime  int64
	ticker       *time.Ticker
	stopCh       chan bool
	wg           *sync.WaitGroup
	running      bool
	mutex        *sync.Mutex
}

func NewBackgroundCleaner(logger echo.Logger, queueManager *QueueManager, store *Store, tickerDuration time.Duration, jobLifetime int64) *BackgroundCleaner {
	return &BackgroundCleaner{
		logger:       logger,
		queueManager: queueManager,
		store:        store,
		jobLifetime:  jobLifetime,
		ticker:       time.NewTicker(tickerDuration),
		stopCh:       make(chan bool),
		wg:           &sync.WaitGroup{},
		running:      false,
		mutex:        &sync.Mutex{},
	}
}

func (bg *BackgroundCleaner) Start() {
	bg.wg.Add(1)
	go func() {
		defer bg.wg.Done()
		for {
			select {
			case <-bg.ticker.C:
				bg.run()
			case <-bg.stopCh:
				return
			}
		}
	}()
}

func (bg *BackgroundCleaner) Stop() {
	bg.ticker.Stop()
	close(bg.stopCh)
	bg.wg.Wait()
}

func (bg *BackgroundCleaner) run() {
	defer func() {
		if r := recover(); r != nil {
			bg.logger.Errorf("BackgroundCleaner caused error: %v", r)
		}
	}()

	bg.logger.Debug("Run the BackgroundCleaner task")

	if !bg.shouldRun() {
		bg.logger.Warn("BackgroundCleaner task has been already running. skip it.")
		return
	}
	defer bg.done()

	tt := time.Now().Add(time.Duration(-1*bg.jobLifetime) * time.Second)
	begin := katsubushi.ToID(tt)
	query := &ListJobsQuery{
		Reverse: true,
		Begin:   &begin,
	}

	bg.logger.Debugf("Try to get before %v (%d) jobs to delete (keep %d sec)", tt, begin, bg.jobLifetime)
	list, err := bg.store.ListJobs(query)
	if err != nil {
		bg.logger.Error(err)
	}
	bg.logger.Debugf("Got %d jobs to delete", list.Count)

	for _, job := range list.Jobs {
		// delete
		if job.Running {
			bg.logger.Debugf("job %d is running. skip it", job.ID)
			continue
		}

		if job.Waiting {
			bg.logger.Debugf("job %d is waiting. skip it", job.ID)
			continue
		}

		if job.FinishedAt == nil {
			bg.logger.Debugf("job %d is not finished. skip it", job.ID)
			continue
		}

		if err := bg.store.DeleteJob(job.ID); err != nil {
			bg.logger.Error(err)
		}
		bg.logger.Debugf("deleted job: %d", job.ID)
	}
}

func (bg *BackgroundCleaner) shouldRun() bool {
	bg.mutex.Lock()
	defer bg.mutex.Unlock()

	if bg.running {
		// already running. should not run.
		return false
	}

	bg.running = true
	return true
}

func (bg *BackgroundCleaner) done() {
	bg.mutex.Lock()
	defer bg.mutex.Unlock()

	bg.running = false
}
