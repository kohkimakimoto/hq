package server

import (
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/hq/structs"
	"github.com/robfig/cron"
	"sync"
	"time"
)

type Background struct {
	App  *App
	Cron *cron.Cron
}

func NewBackground(app *App) *Background {
	return &Background{
		App:  app,
		Cron: cron.New(),
	}
}

func (bg *Background) Start() {
	logger := bg.App.Logger
	logger.Debug("Starting background.")

	app := bg.App
	// bg.Cron.AddFunc("* * * * * *", cleanupJobs(app))
	bg.Cron.AddFunc("@hourly", cleanupJobs(app))
	bg.Cron.Start()
}

func (bg *Background) Close() {
	if bg.Cron != nil {
		bg.Cron.Stop()
	}
}

func cleanupJobs(app *App) func() {
	mutex := new(sync.Mutex)
	logger := app.Logger
	config := app.Config

	return func() {
		if config.KeepJobs <= 0 {
			return
		}

		mutex.Lock()
		defer mutex.Unlock()

		logger.Debug("Run the background task to clean up jobs")

		tt := time.Now().Add(time.Duration(-1*config.KeepJobs) * time.Second)
		begin := katsubushi.ToID(tt)
		query := &structs.ListJobsQuery{
			Reverse:  true,
			Begin:    begin,
			HasBegin: true,
		}

		list := &structs.JobList{
			Jobs:    []*structs.Job{},
			HasNext: false,
		}

		logger.Debugf("Try to get before %v (%d) jobs to delete (keep %d sec)", tt, begin, config.KeepJobs)

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
