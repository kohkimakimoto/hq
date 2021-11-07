package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/kohkimakimoto/hq/internal/command"
	"github.com/kohkimakimoto/hq/internal/version"
)

func main() {
	log.SetFlags(0)
	handleError := func(err interface{}) {
		log.Fatal("error: ", err)
	}

	defer func() {
		if err := recover(); err != nil {
			handleError(err)
		}
	}()

	app := cli.NewApp()
	app.Name = version.Name
	app.Version = version.Version + " (" + version.CommitHash + ")"
	app.Usage = "Simplistic job queue engine"
	app.Copyright = "Copyright (c) 2019 Kohki Makimoto"
	app.Commands = command.Commands
	if err := app.Run(os.Args); err != nil {
		handleError(err)
	}
}
