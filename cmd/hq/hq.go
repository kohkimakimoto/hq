package main

import (
	"github.com/kohkimakimoto/hq/internal/command"
	"github.com/kohkimakimoto/hq/internal/version"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
		}
	}()

	app := cli.NewApp()
	app.Name = version.Name
	app.Version = version.Version + " (" + version.CommitHash + ")"
	app.Usage = "Simplistic job queue engine"
	app.Commands = command.Commands
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
