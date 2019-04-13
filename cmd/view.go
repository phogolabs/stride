package cmd

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/rest/middleware"
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
	router := chi.NewRouter()

	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.NoCache)
	router.Use(middleware.Logger)
	router.Use(middleware.LiveReloader)

	router.Mount("/", http.FileServer(parcello.ManagerAt("viewer")))
	router.Mount("/swagger.spec", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, ctx.String("file-path"))
	}))

	log.Infof("http server is listening on http://%v", ctx.String("listen-addr"))

	return http.ListenAndServe(ctx.String("listen-addr"), router)
}
