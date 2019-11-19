package service

import (
	"io"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/phogolabs/log"
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
	logger := log.GetContext(r.Context())

	file, err := os.Create(e.Path)
	if err != nil {
		logger.WithError(err).Error("failed to create the spec file")
	} else if _, err = io.Copy(file, r.Body); err != nil {
		logger.WithError(err).Error("failed to save the spec fiel")
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// close the file
	file.Close()
}

func (e *Editor) serve(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(parcello.ManagerAt("editor"))
	handler.ServeHTTP(w, r)
}
