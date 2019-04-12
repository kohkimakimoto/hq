package server

import (
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/robfig/cron"
	"sync"
	"time"
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

func (bg *Background) Start() {
	logger := bg.app.Logger
	logger.Debug("Starting background.")

	app := bg.app
	// bg.cron.AddFunc("* * * * * *", cleanupJobs(app))
	bg.cron.AddFunc("@hourly", cleanupJobs(app))
	bg.cron.Start()
}

func (bg *Background) Close() {
	if bg.cron != nil {
		bg.cron.Stop()
	}
}

func cleanupJobs(app *App) func() {
	mutex := new(sync.Mutex)
	logger := app.Logger
	config := app.Config

	return func() {
		if config.JobLifetime <= 0 {
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		logger.Debug("Run the background task to clean up jobs")

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
			if job.FinishedAt != nil {
				if err := app.Store.DeleteJob(job.ID); err != nil {
					logger.Error(err)
				}
				logger.Debugf("deleted job: %d", job.ID)
			}
		}
	}
}
