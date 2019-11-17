package cmd

import (
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/stride/service"
)

// OpenAPIViewer provides a subcommands to view OpenAPI specification in the browser
type OpenAPIViewer struct{}

// CreateCommand creates a cli.Command that can be used by cli.App.
func (m *OpenAPIViewer) CreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "view",
		Usage:       "Shows an OpenAPI specification in the browser",
		Description: "Shows an OpenAPI specification in the browser",
		Before:      m.before,
		Action:      m.view,
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

func (m *OpenAPIViewer) before(ctx *cli.Context) error {
	log.SetHandler(console.New(ctx.Writer))
	return nil
}

func (m *OpenAPIViewer) view(ctx *cli.Context) error {
	var (
		config = &service.ViewerConfig{
			Addr: ctx.String("listen-addr"),
			Path: ctx.String("file-path"),
		}
		server = service.NewViewer(config)
	)

	log.Infof("http server is listening on http://%v", config.Addr)
	return server.ListenAndServe()
}
