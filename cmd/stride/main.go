package main

import (
	"os"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/stride/cmd"
	_ "github.com/phogolabs/stride/template"
)

func main() {
	var (
		editor    = &cmd.OpenAPIEditor{}
		viewer    = &cmd.OpenAPIViewer{}
		generator = &cmd.OpenAPIGenerator{}
		validator = &cmd.OpenAPIValidator{}
		mocker    = &cmd.OpenAPIMocker{}
	)

	commands := []*cli.Command{
		editor.CreateCommand(),
		viewer.CreateCommand(),
		mocker.CreateCommand(),
		generator.CreateCommand(),
		validator.CreateCommand(),
	}

	app := &cli.App{
		Name:      "stride",
		HelpName:  "stride",
		Usage:     "OpenAPI viewer, editor, generator and validator",
		UsageText: "stride [global options]",
		Version:   "1.0-beta-05",
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		Commands:  commands,
	}

	app.Run(os.Args)
}
