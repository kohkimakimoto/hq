package server

import (
	"context"
	"github.com/boltdb/bolt"
	"github.com/kayac/go-katsubushi"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/kohkimakimoto/hq/util/logutil"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"html/template"
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
	// LogFile is log file is used by Logger.
	Logfile *os.File
	// LogLevel
	LogLevel log.Lvl
	// Echo web framework
	Echo *echo.Echo
	// AccessLog
	AccessLogFile *os.File
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
	// View
	View *template.Template
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
		Config:        c,
		Echo:          echo.New(),
		AccessLogFile: os.Stdout,
		DataDir:       c.DataDir,
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
		gen:    gen,
	}

	if err := app.Store.Init(); err != nil {
		return err
	}

	// background
	app.Background = NewBackground(app)

	return nil
}

func (app *App) openLogfile() error {
	if app.Config.Logfile != "" {
		f, err := os.OpenFile(app.Config.Logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return err
		}

		app.Logger.SetOutput(f)
		app.Logfile = f
	}

	return nil
}

func (app *App) BoltDBPath() string {
	return filepath.Join(app.DataDir, "server.bolt")
}

func (app *App) ListenAndServe() error {
	e := app.Echo

	go func() {
		if err := e.Start(app.Config.Addr); err != nil {
			e.Logger.Info(err)
		}
	}()

	app.Logger.Infof("The server Listening on %s (pid: %d)", e.Server.Addr, os.Getpid())

	// start background process
	app.Background.Start()

	// see https://echo.labstack.com/cookbook/graceful-shutdown
	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	e.Logger.Info("Shutting down the server")

	if err := e.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "fail to shutdown")
	}

	return nil
}

func (app *App) Close() error {
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
