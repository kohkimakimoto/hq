package server

import (
	"github.com/robfig/cron"
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

	//srv := bg.Server
	//bg.Cron.AddFunc("* * * * * *", cleanupResources(srv))
	// bg.Cron.AddFunc("@hourly", cleanupResources(srv))

	bg.Cron.Start()
}
