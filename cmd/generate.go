package cmd

import (
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
)

// OpenAPIGenerator provides a subcommands to generate source code from OpenAPI specification
type OpenAPIGenerator struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *OpenAPIGenerator) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "generate",
		Usage:       "Generates a project from an OpenAPI",
		Description: "Generates a project from an OpenAPI",
		Before:      m.before,
		Action:      m.generate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file-path",
				Usage: "path to the open api specification",
				Value: "./open-api.yaml",
			},
		},
	}
}

func (m *OpenAPIGenerator) before(ctx *cli.Context) error {
	log.SetHandler(console.New(ctx.Writer))
	return nil
}

func (m *OpenAPIGenerator) generate(ctx *cli.Context) error {
	return nil
}
