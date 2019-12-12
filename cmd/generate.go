package cmd

import (
	"path/filepath"

	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/stride/codedom"
	"github.com/phogolabs/stride/service"
	"github.com/phogolabs/stride/syntax/golang"
	"github.com/phogolabs/stride/syntax/markdown"
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
				Value: "./swagger.yaml",
			},
			&cli.StringFlag{
				Name:   "project-path, p",
				Usage:  "path to the project directory",
				Value:  ".",
				EnvVar: "PWD",
			},
		},
	}
}

func (m *OpenAPIGenerator) before(ctx *cli.Context) error {
	log.SetHandler(console.New(ctx.Writer))
	return nil
}

func (m *OpenAPIGenerator) generate(ctx *cli.Context) error {
	// get the spec
	path, err := get(ctx, "file-path")
	if err != nil {
		return err
	}

	dir, err := filepath.Abs(ctx.String("project-path"))
	if err != nil {
		return err
	}

	// generate the soec
	generator := &service.Generator{
		Path: path,
		Resolver: &codedom.Resolver{
			Reporter: reporter(ctx),
			Cache:    codedom.TypeDescriptorMap{},
		},
		Generator: service.CompositeGenerator{
			&golang.Generator{
				Reporter: reporter(ctx),
				Path:     dir,
			},
			&markdown.Generator{
				Reporter: reporter(ctx),
				Path:     dir,
			},
		},
	}

	return generator.Generate()
}
