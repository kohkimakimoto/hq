package command

import (
	"github.com/kohkimakimoto/hq/server"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pkg/errors"
	"github.com/urfave/cli"

	"os"
)

var ServeCommand = cli.Command{
	Name:   "serve",
	Usage:  "Start the hq server process",
	Action: serverAction,
	Flags: []cli.Flag{
		configFileFlag,
		logLevelFlag,
	},
}

func serverAction(ctx *cli.Context) error {
	config := server.NewConfig()

	if err := loadServerConfigFiles(ctx, config); err != nil {
		return err
	}

	applyLogLevel(ctx, config)

	app := server.NewApp(config)
	defer app.Close()

	if err := configureServer(app); err != nil {
		return errors.Wrapf(err, "failed to initialize hq server")
	}

	return app.ListenAndServe()
}

func loadServerConfigFiles(ctx *cli.Context, config *server.Config) error {
	if v := ctx.String("config-file"); v != "" {
		// config file specified by CLI option
		if err := config.LoadConfigFile(v); err != nil {
			return errors.Wrapf(err, "failed to load config from the file '%s'", v)
		}
	} else {
		// default system config
		configFile := "/etc/hq/hq.toml"

		if _, err := os.Stat(configFile); err == nil {
			if err := config.LoadConfigFile(configFile); err != nil {
				return errors.Wrapf(err, "failed to load config from the file '%s'", configFile)
			}
		}
	}

	return nil
}

func configureServer(app *server.App) error {
	// open resources such as database, temporary directory, etc.
	if err := app.Open(); err != nil {
		return err
	}

	e := app.Echo
	c := app.Config

	// open access log file
	accessLogfile := os.Stdout
	if c.AccessLogfile != "" {
		f, err := os.OpenFile(c.AccessLogfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return errors.Wrapf(err, "failed to open access logfile")
		}
		accessLogfile = f
	}

	// error handler
	e.HTTPErrorHandler = errorHandler(app)

	// middleware
	e.Use(server.AppContextMiddleware(app))
	e.Use(middleware.Recover())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format:  `${time_rfc3339} ${remote_ip} ${host} ${method} ${uri} ${status} ${latency} ${latency_human} ${bytes_in} ${bytes_out}` + "\n",
		Output:  accessLogfile,
	}))

	// handlers
	e.Any("/", server.InfoHandler)
	e.POST("/job", server.CreateJobHandler)
	e.GET("/job", server.ListJobsHandler)
	e.GET("/job/:id", server.GetJobHandler)
	e.POST("/job/:id/restart", server.RestartJobHandler)
	e.POST("/job/:id/stop", server.StopJobHandler)
	e.GET("/stats", server.StatsHandler)
	e.DELETE("/job/:id", server.DeleteJobHandler)

	return nil
}

func errorHandler(app *server.App) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		e := c.Echo()
		if c.Response().Committed {
			goto ERROR
		}

		server.ErrorHandler(err, c)
	ERROR:
		e.Logger.Error(err)
	}
}
