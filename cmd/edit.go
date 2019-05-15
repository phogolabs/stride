package cmd

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/phogolabs/cli"
	"github.com/phogolabs/log"
	"github.com/phogolabs/log/handler/console"
	"github.com/phogolabs/parcello"
	"github.com/phogolabs/rest/middleware"
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
	router := chi.NewRouter()

	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.NoCache)
	router.Use(middleware.Logger)

	router.Mount("/", http.FileServer(parcello.ManagerAt("editor")))
	// router.Mount("/", http.FileServer(http.Dir("./template/editor")))
	router.Mount("/swagger.spec", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.ServeFile(w, r, ctx.String("file-path"))
		case "POST":
			SaveFile(w, r, ctx.String("file-path"))
		}
	}))

	log.Infof("http server is listening on http://%v", ctx.String("listen-addr"))

	return http.ListenAndServe(ctx.String("listen-addr"), router)
}

// SaveFile saves the file
func SaveFile(w http.ResponseWriter, r *http.Request, path string) {
	spec, err := os.Create(path)
	if err != nil {
		middleware.GetLogger(r).WithError(err).Error("failed to save the spec")
		return
	}
	defer spec.Close()

	if _, err := io.Copy(spec, r.Body); err != nil {
		middleware.GetLogger(r).WithError(err).Error("failed to save the spec")
	}
}
