package command

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/cheynewallace/tabby"
	"github.com/urfave/cli/v2"

	"github.com/kohkimakimoto/hq/internal/client"
)

var Commands = []*cli.Command{
	DeleteCommand,
	InfoCommand,
	ListCommand,
	PushCommand,
	RestartCommand,
	ServeCommand,
	StatsCommand,
	StopCommand,
}

// Flags
var (
	logLevelFlag = &cli.StringFlag{
		Name:    "log-level",
		Aliases: []string{"l"},
		Usage:   "Set Log `LEVEL` (error|warn|info|debug).",
		EnvVars: []string{"HQ_LOG_LEVEL"},
	}
	configFileFlag = &cli.StringFlag{
		Name:    "config-file",
		Aliases: []string{"c"},
		Usage:   "Load config from the `FILE`",
		EnvVars: []string{"HQ_CONFIG"},
	}
	addressFlag = &cli.StringFlag{
		Name:    "address",
		Aliases: []string{"a"},
		Usage:   "The `ADDRESS` of the HQ server.",
		Value:   "http://127.0.0.1:19900",
	}
)

type ClientFactory func(c *cli.Context) *client.Client

func defaultClientFactory(c *cli.Context) *client.Client {
	return client.New(c.String("address"))
}

func setClientFactory(app *cli.App, factory ClientFactory) {
	if app.Metadata == nil {
		app.Metadata = make(map[string]interface{})
	}

	app.Metadata["clientFactory"] = factory
}

func newClient(ctx *cli.Context) *client.Client {
	var factory ClientFactory
	f, ok := ctx.App.Metadata["clientFactory"]
	if !ok {
		factory = defaultClientFactory
	} else {
		ff, ok := f.(ClientFactory)
		if !ok {
			panic(fmt.Sprintf("invalid client factory %v", f))
		}
		factory = ff
	}

	return factory(ctx)
}

func newTabby(output io.Writer) *tabby.Tabby {
	return tabby.NewCustom(tabwriter.NewWriter(output, 0, 0, 2, ' ', 0))
}
