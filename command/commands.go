package command

import (
	"os"
	"runtime"
	"text/tabwriter"

	"github.com/cheynewallace/tabby"
	"github.com/kohkimakimoto/hq/client"
	"github.com/urfave/cli"
)

// Command set
var Commands = []cli.Command{
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
	// common flags
	logLevelFlag = cli.StringFlag{
		Name:   "log-level, l",
		Usage:  "Set Log `LEVEL` (error|warn|info|debug).",
		EnvVar: "HQ_LOG_LEVEL",
	}

	// serve flags
	configFileFlag = cli.StringFlag{
		Name:   "config-file, c",
		Usage:  "Load config from the `FILE`",
		EnvVar: "HQ_CONFIG",
	}

	// client flag
	addressFlag = cli.StringFlag{
		Name:  "address, a",
		Usage: "The `ADDRESS` of the HQ server.",
		Value: "http://127.0.0.1:19900",
	}
)

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

type LogLevelSetter interface {
	SetLogLevel(string)
}

func applyLogLevel(ctx *cli.Context, setter LogLevelSetter) {
	if v := ctx.String("log-level"); v != "" {
		setter.SetLogLevel(v)
	}
}

func newClient(ctx *cli.Context) *client.Client {
	return client.New(ctx.String("address"))
}

func newTabby() *tabby.Tabby {
	return tabby.NewCustom(tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0))
}

func init() {
	cli.AppHelpTemplate = `Usage: {{.HelpName}}{{if .VisibleFlags}} [<options...>]{{end}} <command>

{{if .Usage}}{{.Usage}}{{end}}{{if .Version}}
version {{.Version}}{{end}}{{if .Flags}}

Options:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}{{if .VisibleCommands}}
Commands:{{range .VisibleCategories}}{{if .Name}}
  {{.Name}}:{{end}}{{range .VisibleCommands}}
  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{end}}
`
	cli.CommandHelpTemplate = `Usage: {{.HelpName}}{{if .VisibleFlags}} [<options...>]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{end}}
{{if .Description}}
{{.Description}}
{{end -}}
{{if .VisibleFlags}}
Options:
  {{range .VisibleFlags}}{{.}}
  {{end}}{{end}}
`

	cli.SubcommandHelpTemplate = `Usage: {{.HelpName}} <subcommand>
{{if .Description}}
{{.Description}}
{{end -}}
{{if .VisibleCommands}}
Subcommands:{{range .VisibleCategories}}{{if .Name}}
  {{.Name}}:{{end}}{{range .VisibleCommands}}
  {{join .Names ", "}}{{"\t"}}{{.Usage}}{{end}}
{{end}}{{end}}
`
}
