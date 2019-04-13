package cmd

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/stride/codegen"
)

// OpenAPIGenerator provides a subcommands to generate source code from OpenAPI specification
type OpenAPIGenerator struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *OpenAPIGenerator) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "generate",
		Usage:       "Generates a project from an OpenAPI specification",
		Description: "Generates a project from an OpenAPI specification",
		Before:      m.before,
		Action:      m.generate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file-path, f",
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
	var (
		loader = openapi3.NewSwaggerLoader()
		path   = ctx.String("file-path")
	)

	spec, err := loader.LoadSwaggerFromFile(path)
	if err != nil {
		return err
	}

	resolver := &codegen.Resolver{
		Schemas: make(map[string]*codegen.TypeDescriptor),
	}

	resolver.Resolve(spec)
	return nil
}
