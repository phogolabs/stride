package cmd

import (
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/stride/service"
)

// OpenAPIValidator provides a subcommands to view OpenAPI specification in the browser
type OpenAPIValidator struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *OpenAPIValidator) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "validate",
		Usage:       "Validates an OpenAPI specification",
		Description: "Validates an OpenAPI specification",
		Before:      m.before,
		Action:      m.validate,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file-path, f",
				Usage: "path to the open api specification",
				Value: "./swagger.yaml",
			},
		},
	}
}

func (m *OpenAPIValidator) before(ctx *cli.Context) error {
	log.SetHandler(console.New(ctx.Writer))
	return nil
}

func (m *OpenAPIValidator) validate(ctx *cli.Context) error {
	// get the spec
	path, err := get(ctx, "file-path")
	if err != nil {
		return err
	}

	validator := &service.Validator{
		Path:     path,
		Reporter: reporter(ctx),
	}

	return validator.Validate()
}
