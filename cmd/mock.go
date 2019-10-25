package cmd

import (
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
)

// OpenAPIMocker provides a subcommands to generate source code from OpenAPI specification
type OpenAPIMocker struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *OpenAPIMocker) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "mock",
		Usage:       "Runs a mock server from an OpenAPI specification",
		Description: "Runs a mock server from an OpenAPI specification",
		Before:      m.before,
		Action:      m.mock,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file-path, f",
				Usage: "path to the open api specification",
				Value: "./swagger.yaml",
			},
		},
	}
}

func (m *OpenAPIMocker) before(ctx *cli.Context) error {
	log.SetHandler(console.New(ctx.Writer))
	return nil
}

func (m *OpenAPIMocker) mock(ctx *cli.Context) error {
	return nil
}
