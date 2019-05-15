package main

import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/stride/cmd"
	_ "github.com/phogolabs/stride/template"
)

func main() {
	var (
		generator = &cmd.OpenAPIGenerator{}
		viewer    = &cmd.OpenAPIViewer{}
		editor    = &cmd.OpenAPIEditor{}
	)

	commands := []*cli.Command{
		generator.CreateCommand(),
		viewer.CreateCommand(),
		editor.CreateCommand(),
	}

	app := &cli.App{
		Name:      "stride",
		HelpName:  "stride",
		Usage:     "OpenAPI Viewer and Generator",
		UsageText: "stride [global options]",
		Version:   "1.0-beta-05",
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Commands:  commands,
	}

	if err := app.Run(os.Args); err != nil {
		log.WithError(err).Fatal("execution failed")
	}
}
