package server

import (
	"sync"
	"time"

	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/robfig/cron/v3"
)

type Background struct {
	app  *App
	cron *cron.Cron
}

func NewBackground(app *App) *Background {
	return &Background{
		app:  app,
		cron: cron.New(),
	}
}

func (bg *Background) Start() error {
	logger := bg.app.Logger
	logger.Debug("Starting background.")

	config := bg.app.Config
	app := bg.app

	if config.JobLifetime > 0 {
		if _, err := bg.cron.AddFunc("@hourly", cleanupJobs(app)); err != nil {
			return err
		}
	}

	bg.cron.Start()

	return nil
}

func (bg *Background) Close() {
	if bg.cron != nil {
		bg.cron.Stop()
	}
}

func cleanupJobs(app *App) func() {
	logger := app.Logger
	config := app.Config

	running := false
	mutex := new(sync.Mutex)

	return func() {
		logger.Debug("Run the background task 'cleanupJobs'")
		if running {
			logger.Warn("'cleanupJobs' has been already running. skip it.")
			return
		}

		mutex.Lock()
		running = true
		mutex.Unlock()

		defer func() {
			running = false
		}()

		tt := time.Now().Add(time.Duration(-1*config.JobLifetime) * time.Second)
		begin := katsubushi.ToID(tt)
		query := &ListJobsQuery{
			Reverse: true,
			Begin:   &begin,
		}

		list := &hq.JobList{
			Jobs:    []*hq.Job{},
			HasNext: false,
		}

		logger.Debugf("Try to get before %v (%d) jobs to delete (keep %d sec)", tt, begin, config.JobLifetime)

		if err := app.Store.ListJobs(query, list); err != nil {
			logger.Error(err)
		}

		logger.Debugf("Got %d jobs to delete", list.Count)

		for _, job := range list.Jobs {
			// delete
			if job.Running {
				logger.Debugf("job %d is running. skip it", job.ID)
				continue
			}

			if job.Waiting {
				logger.Debugf("job %d is waiting. skip it", job.ID)
				continue
			}

			if job.FinishedAt == nil {
				logger.Debugf("job %d is not finished. skip it", job.ID)
				continue
			}

			if err := app.Store.DeleteJob(job.ID); err != nil {
				logger.Error(err)
			}
			logger.Debugf("deleted job: %d", job.ID)

		}
	}
}
