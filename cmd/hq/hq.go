package main

import (
	"fmt"
	"github.com/kohkimakimoto/hq/hq"
	"github.com/kohkimakimoto/hq/command"
	"os"
	"github.com/urfave/cli"
)

func main() {
	os.Exit(realMain())
}

func realMain() (status int) {
	defer func() {
		if err := recover(); err != nil {
			printError(err)
			status = 1
		}
	}()

	app := cli.NewApp()
	app.Name = hq.Name
	app.HelpName = hq.DisplayName
	app.Version = hq.Version + " (" + hq.CommitHash + ")"
	app.Usage = "Server Management Framework"
	app.Commands = command.Commands

	if err := app.Run(os.Args); err != nil {
		printError(err)
		status = 1
	}

	return status
}

func printError(err interface{}) {
	fmt.Fprintf(os.Stderr, "%s %v\n", "ERROR", err)
}

