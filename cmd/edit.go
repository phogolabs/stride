package cmd

import (
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/stride/service"
)

// OpenAPIEditor provides a subcommands to edit OpenAPI specification in the browser
type OpenAPIEditor struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *OpenAPIEditor) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "edit",
		Usage:       "Edit an OpenAPI specification in the browser",
		Description: "Edit an OpenAPI specification in the browser",
		Before:      m.before,
		Action:      m.edit,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen-addr",
				Usage: "address on which the http server is listening on",
				Value: "localhost:8080",
			},
			&cli.StringFlag{
				Name:  "file-path, f",
				Usage: "path to the open api specification",
				Value: "./swagger.yaml",
			},
		},
	}
}

func (m *OpenAPIEditor) before(ctx *cli.Context) error {
	log.SetHandler(console.New(ctx.Writer))
	return nil
}

func (m *OpenAPIEditor) edit(ctx *cli.Context) error {
	var (
		config = &service.EditorConfig{
			Addr: ctx.String("listen-addr"),
			Path: ctx.String("file-path"),
		}
		server = service.NewEditor(config)
	)

	log.Infof("http server is listening on http://%v", config.Addr)
	return server.ListenAndServe()
}
