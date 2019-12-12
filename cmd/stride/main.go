package main

import (
	"os"
	"os/signal"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/stride/async"
	"github.com/phogolabs/stride/cmd"

	_ "github.com/phogolabs/stride/template"
)

func main() {
	var (
		editor    = &cmd.OpenAPIEditor{}
		viewer    = &cmd.OpenAPIViewer{}
		generator = &cmd.OpenAPIGenerator{}
		validator = &cmd.OpenAPIValidator{}
		// mocker    = &cmd.OpenAPIMocker{}
	)

	commands := []*cli.Command{
		editor.CreateCommand(),
		viewer.CreateCommand(),
		// mocker.CreateCommand(),
		generator.CreateCommand(),
		validator.CreateCommand(),
	}

	app := &cli.App{
		Name:      "stride",
		HelpName:  "stride",
		Usage:     "OpenAPI viewer, editor, generator, validator and mocker",
		UsageText: "stride [global options]",
		Version:   "1.0-beta-05",
		Writer:    os.Stdout,
		ErrWriter: os.Stderr,
		OnSignal:  onSignal,
		Commands:  commands,
	}

	app.Run(os.Args)
}

func onSignal(ctx *cli.Context, term os.Signal) error {
	if term == os.Interrupt {
		if task, ok := ctx.Metadata["task"].(*async.Task); ok {
			signal.Reset(term)
			return task.Stop()
		}
	}

	return nil
}
