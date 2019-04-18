package server

import (
	"context"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/client9/reopen"
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/kohkimakimoto/hq/util/logutil"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
)

type App struct {
	// Configuration of the application instance
	Config *Config
	// Logger
	Logger echo.Logger
	// LogfileWriter
	LogfileWriter reopen.Writer
	// LogLevel
	LogLevel log.Lvl
	// Echo web framework
	Echo *echo.Echo
	// AccessLog
	AccessLogFile *os.File
	// AccessLogFile
	AccessLogFileWriter reopen.Writer
	// DataDir
	DataDir string
	// UseTempDataDir
	UseTempDataDir bool
	// DB
	DB *bolt.DB
	// Store
	Store *Store
	// Background
	Background *Background
	// katsubushi
	Gen katsubushi.Generator
	// QueueManager
	QueueManager *QueueManager
}

func NewApp(config ...*Config) *App {
	var c *Config
	if len(config) == 0 {
		c = NewConfig()
	} else {
		c = config[0]
	}

	// create app instance
	app := &App{
		Config:  c,
		Echo:    echo.New(),
		DataDir: c.DataDir,
	}

	app.Echo.HideBanner = true
	app.Echo.HidePort = true
	app.Echo.Server.Addr = app.Config.Addr

	return app
}

func (app *App) Open() error {
	config := app.Config

	// log level
	lv, err := logutil.LoglvlFromString(config.LogLevel)
	if err != nil {
		return err
	}
	app.LogLevel = lv

	// setup logger
	logger := log.New(hq.Name)
	logger.SetLevel(app.LogLevel)
	logger.SetHeader(`${time_rfc3339} ${level}`)
	app.Logger = logger
	// setup echo logger
	app.Echo.Logger = logger

	// open log
	if err := app.openLogfile(); err != nil {
		return err
	}

	// Uniqid generator
	epoch, err := config.IDEpochTime()
	if err != nil {
		return err
	}
	katsubushi.Epoch = epoch
	gen, err := katsubushi.NewGenerator(config.ServerId)
	if err != nil {
		return err
	}
	app.Gen = gen

	// setup data directory as a temporary directory if it is not set.
	if app.DataDir == "" {
		logger.Warn("Your 'data_dir' configuration is not set. It uses a temporary directory that is deleted after the process terminates.")

		tmpdir, err := ioutil.TempDir("", hq.Name+"_data_")
		if err != nil {
			return err
		}
		logger.Warnf("Created temporary data directory: %s", tmpdir)
		app.DataDir = tmpdir
		app.UseTempDataDir = true
	}

	if _, err := os.Stat(app.DataDir); os.IsNotExist(err) {
		err = os.MkdirAll(app.DataDir, os.FileMode(0755))
		if err != nil {
			return err
		}
	}

	logger.Infof("Opened data directory: %s", app.DataDir)

	// setup bolt database
	db, err := bolt.Open(app.BoltDBPath(), 0600, nil)
	if err != nil {
		return err
	}
	app.DB = db
	logger.Infof("Opened boltdb: %s", db.Path())

	// store
	app.Store = &Store{
		app:    app,
		db:     db,
		logger: logger,
	}

	if err := app.Store.Init(); err != nil {
		return err
	}

	// queue
	app.QueueManager = NewQueueManager(app)
	app.QueueManager.Start()

	// background
	app.Background = NewBackground(app)
	app.Background.Start()

	return nil
}

func (app *App) openLogfile() error {
	if app.Config.Logfile != "" {
		f, err := reopen.NewFileWriterMode(app.Config.Logfile, 0644)
		if err != nil {
			return err
		}

		app.Logger.SetOutput(f)
		app.LogfileWriter = f
	} else {
		app.LogfileWriter = reopen.Stdout
	}

	if app.Config.AccessLogfile != "" {
		f, err := reopen.NewFileWriterMode(app.Config.AccessLogfile, 0644)
		if err != nil {
			return err
		}
		app.AccessLogFileWriter = f
	} else {
		app.AccessLogFileWriter = reopen.Stdout
	}

	return nil
}

func (app *App) BoltDBPath() string {
	return filepath.Join(app.DataDir, "server.bolt")
}

func (app *App) ListenAndServe() error {
	// open resources such as log files, database, temporary directory, etc.
	if err := app.Open(); err != nil {
		return err
	}

	// Configure http servers (handlers and middleware)
	e := app.Echo

	// error handler
	e.HTTPErrorHandler = errorHandler(app)
	// middleware
	e.Use(AppContextMiddleware(app))
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  `${time_rfc3339} ${remote_ip} ${host} ${method} ${uri} ${status} ${latency} ${latency_human} ${bytes_in} ${bytes_out}` + "\n",
		Output:  app.AccessLogFileWriter,
	}))
	// handlers
	e.Any("/", InfoHandler)
	e.POST("/job", CreateJobHandler)
	e.GET("/job", ListJobsHandler)
	e.GET("/job/:id", GetJobHandler)
	e.POST("/job/:id/restart", RestartJobHandler)
	e.POST("/job/:id/stop", StopJobHandler)
	e.GET("/stats", StatsHandler)
	e.DELETE("/job/:id", DeleteJobHandler)

	// handler for reopen logs
	go app.sigusr1Handler()

	// start server.
	go func() {
		if err := e.Start(app.Config.Addr); err != nil {
			e.Logger.Info(err)
		}
	}()

	app.Logger.Infof("The server Listening on %s (pid: %d)", e.Server.Addr, os.Getpid())

	// see https://echo.labstack.com/cookbook/graceful-shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	app.Logger.Infof("Received signal: %v", sig)
	timeout := time.Duration(app.Config.ShutdownTimeout) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	e.Logger.Info("Shutting down the server")

	if err := e.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "fail to shutdown echo http server")
	}

	// wait for running jobs finishing.
	e.Logger.Info("Waiting for finishing the jobs.")
	app.QueueManager.Wait()
	app.Logger.Infof("Successfully shutdown")

	return nil
}

func (app *App) sigusr1Handler() {
	reopen := make(chan os.Signal, 1)
	signal.Notify(reopen, syscall.SIGUSR1)

	logger := app.Logger

	for {
		select {
		case sig := <-reopen:
			logger.Infof("Received signal to reopen the logs: %v", sig)

			if err := app.LogfileWriter.Reopen(); err != nil {
				logger.Error(fmt.Sprintf("failed to reopen log: %v", err))
			}

			if err := app.AccessLogFileWriter.Reopen(); err != nil {
				logger.Error(fmt.Sprintf("failed to reopen access log: %v", err))
			}
		}
	}
}

func (app *App) Close() error {
	if app.Background != nil {
		app.Background.Close()
	}

	if app.DB != nil {
		if err := app.DB.Close(); err != nil {
			return err
		}
	}

	if app.UseTempDataDir {
		if err := os.RemoveAll(app.DataDir); err != nil {
			return err
		}
		app.Logger.Warnf("Removed temporary directory: %s", app.DataDir)
	}

	return nil
}

type AppContext struct {
	echo.Context
	app *App
}

func (c *AppContext) App() *App {
	return c.app
}

func AppContextMiddleware(app *App) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &AppContext{c, app}
			return next(cc)
		}
	}
}
