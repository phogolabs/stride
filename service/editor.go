package service

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/phogolabs/parcello"
)

// EditorConfig represents the editor config
type EditorConfig struct {
	Addr string
	Path string
}

// NewEditor creates a new editor
func NewEditor(config *EditorConfig) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.NoCache)
	router.Use(middleware.Logger)

	handler := &Editor{
		Path: config.Path,
	}

	handler.Mount(router)

	return &http.Server{
		Addr:    config.Addr,
		Handler: router,
	}
}

// Editor edits the swagger file
type Editor struct {
	Path string
}

// Mount mounts the editor
func (e *Editor) Mount(r chi.Router) {
	r.Get("/*", e.serve)
	r.Get("/swagger.spec", e.load)
	r.Post("/swagger.spec", e.save)
}

func (e *Editor) load(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, e.Path)
}

func (e *Editor) save(w http.ResponseWriter, r *http.Request) {
	spec, err := os.Create(e.Path)
	if err != nil {
		// middleware.GetLogger(r).WithError(err).Error("failed to save the spec")
		return
	}
	defer spec.Close()

	if _, err := io.Copy(spec, r.Body); err != nil {
		// middleware.GetLogger(r).WithError(err).Error("failed to save the spec")
	}
}

func (e *Editor) serve(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(parcello.ManagerAt("editor"))
	handler.ServeHTTP(w, r)
}
