package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/client9/reopen"
	"github.com/kayac/go-katsubushi"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"

	"github.com/kohkimakimoto/hq/ui"
)

func Run(config *Config) error {
	if err := InitApp(config); err != nil {
		return err
	}
	defer g.Close()
	return g.Start()
}

// g is an *App instance that is stored in global space.
// It aims to be used in handlers for easy accessing.
var g *App

// InitApp creates new *App instance and registers it in global space.
func InitApp(config *Config) error {
	app, err := NewApp(config)
	if err != nil {
		return err
	}
	g = app

	return nil
}

// App is an object that represents HQ server.
// It is a global singleton instance.
type App struct {
	// Config is a configuration for the app.
	Config *Config
	// Echo web framework
	Echo *echo.Echo
	// ShutdownTimeoutSec is the timeout for graceful shutdown.
	ShutdownTimeoutSec int64
	// LogFileWriter is the writer for log file.
	LogFileWriter reopen.Writer
	// AccessLogFileWriter is the writer for access log file.
	AccessLogFileWriter reopen.Writer
	// IdGen is an ID generator powered by katsubushi.
	IdGen katsubushi.Generator
	// QueueManager is a Queue manager.
	QueueManager *QueueManager
	// Store is a main database representation.
	Store *Store
	// BackgroundCleaner is a background task runner to clean the stale jobs.
	BackgroundCleaner *BackgroundCleaner
	// Dispatchers
	Dispatchers []*Dispatcher
}

// NewApp creates a new App instance.
func NewApp(c *Config) (*App, error) {
	// echo
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Server.Addr = c.Addr

	// logger
	logger := log.New("hq")
	logger.SetHeader(`${time_rfc3339} ${level}`)
	level, err := c.LogLevel()
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)
	e.Logger = logger

	// server app instance
	a := &App{
		Config:             c,
		Echo:               e,
		ShutdownTimeoutSec: c.ShutdownTimeout,
	}

	// setup log file if it is specified.
	if c.Logfile != "" {
		f, err := reopen.NewFileWriterMode(c.Logfile, 0644)
		if err != nil {
			return nil, err
		}
		a.LogFileWriter = f
		e.Logger.SetOutput(f)
	} else {
		a.LogFileWriter = reopen.Stdout
	}

	// setup access log file if it is specified.
	if c.AccessLogfile != "" {
		f, err := reopen.NewFileWriterMode(c.AccessLogfile, 0644)
		if err != nil {
			return nil, err
		}
		a.AccessLogFileWriter = f
	} else {
		a.AccessLogFileWriter = reopen.Stdout
	}

	// setup ID generator
	epoch, err := c.IDEpochTime()
	if err != nil {
		return nil, err
	}
	katsubushi.Epoch = epoch
	gen, err := katsubushi.NewGenerator(c.ServerId)
	if err != nil {
		return nil, err
	}
	a.IdGen = gen

	// setup Queue manager
	a.QueueManager = NewQueueManager(c.Queues)

	// setup db
	a.Store = NewStore(c.DataDir, e.Logger, a.QueueManager)
	if err := a.Store.Open(); err != nil {
		return nil, err
	}

	// setup background
	a.BackgroundCleaner = NewBackgroundCleaner(e.Logger, a.QueueManager, a.Store, 1*time.Minute, c.JobLifetime)

	// setup dispatchers
	for i := int64(0); i < c.Dispatchers; i++ {
		a.Dispatchers = append(a.Dispatchers, &Dispatcher{
			queueManager:      a.QueueManager,
			store:             a.Store,
			logger:            e.Logger,
			httpClientFactory: defaultHttpClientFactory,
			maxWorkers:        c.MaxWorkers,
			numWorkers:        0,
		})
	}

	// error handler
	e.HTTPErrorHandler = ErrorHandler

	// setup renderer (it is only used for UI)
	e.Renderer = NewUITemplateRenderer(c.UIBasename)

	// middleware
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  `${time_rfc3339} ${remote_ip} ${host} ${method} ${uri} ${status} ${latency} ${latency_human} ${bytes_in} ${bytes_out}` + "\n",
		Output:  a.AccessLogFileWriter,
	}))

	registerAPIHandlers(e, "/")

	if c.UI {
		uiPath := strings.TrimSuffix(c.UIBasename, "/")
		// redirect `/ui` > `/ui/`
		e.Any(uiPath, func(c echo.Context) error {
			return c.Redirect(http.StatusFound, uiPath+"/")
		})
		// enable web ui handlers
		e.Any(uiPath+"/", UIIndexHandler)
		e.Any(uiPath+"/*", UIFallbackHandler)
		e.GET(uiPath+"/dist/*", echo.WrapHandler(http.StripPrefix(uiPath+"/dist/", ui.AssetHandler)))

		// api for web ui
		e.GET(uiPath+"/internal/dashboard", UIDashboardApiHandler)
		// The general api endpoints that are used by web ui
		registerAPIHandlers(e, uiPath+"/internal/")
	}

	return a, nil
}

// Start starts HQ server.
func (a *App) Start() error {
	e := a.Echo
	logger := e.Logger

	// register signal handler
	go a.signalHandler()

	// start dispatchers
	a.startDispatchers()
	logger.Debugf("Started %d dispatcher(s)", len(a.Dispatchers))

	// start background
	a.BackgroundCleaner.Start()
	logger.Debug("Started BackgroundCleaner thread.")

	// start http server
	go func() {
		if err := e.Start(e.Server.Addr); err != nil {
			logger.Info(err)
		}
	}()
	logger.Infof("The server listening on %s (pid: %d)", e.Server.Addr, os.Getpid())

	// see https://echo.labstack.com/cookbook/graceful-shutdown
	// Wait for interrupt signal to gracefully shut down the server with a timeout of 'ShutdownTimeoutSec' seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	logger.Infof("Received signal: %v", sig)
	timeout := time.Duration(a.ShutdownTimeoutSec) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Graceful restart is not supported.
	// Because the HQ server uses boltdb that is locked by the one process.
	// It means that the server can't run multiple processes that is needed for graceful restart.

	logger.Info("Shutting down the server")
	if err := e.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "failed to shut down echo http server")
	}

	// closing workers
	logger.Info("Waiting for finishing the running jobs")
	a.waitDispatchers()
	logger.Info("Finished the jobs")

	// stopping background
	logger.Debug("Stopping BackgroundCleaner")
	a.BackgroundCleaner.Stop()
	logger.Debug("Stopped BackgroundCleaner")

	// done
	logger.Infof("Successfully shutdown")
	return nil
}

func (a *App) startDispatchers() {
	for _, d := range a.Dispatchers {
		go d.EventLoop()
	}
}

func (a *App) waitDispatchers() {
	for _, d := range a.Dispatchers {
		d.Wait()
	}
}

func (a *App) signalHandler() {
	logger := a.Echo.Logger
	reopenSig := make(chan os.Signal, 1)
	// USR1 signal reopens log files
	// It is the same behavior as nginx.
	// see https://www.nginx.com/nginx-wiki/build/dirhtml/start/topics/examples/logrotation/
	signal.Notify(reopenSig, syscall.SIGUSR1)

	for {
		select {
		case sig := <-reopenSig:
			logger.Infof("Received signal to reopen the logs: %v", sig)

			if err := a.LogFileWriter.Reopen(); err != nil {
				logger.Error(fmt.Sprintf("failed to reopen log: %v", err))
			}

			if err := a.AccessLogFileWriter.Reopen(); err != nil {
				logger.Error(fmt.Sprintf("failed to reopen access log: %v", err))
			}
		}
	}
}

func (a *App) Close() {
	if a.Store != nil {
		a.Store.Close()
	}
}
