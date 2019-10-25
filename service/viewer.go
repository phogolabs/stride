package service

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/phogolabs/parcello"
)

// ViewerConfig represents the viewer config
type ViewerConfig struct {
	Addr string
	Path string
}

// NewViewer creates a new viewer
func NewViewer(config *ViewerConfig) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.NoCache)
	router.Use(middleware.Logger)

	handler := &Viewer{
		Path: config.Path,
	}

	handler.Mount(router)

	return &http.Server{
		Addr:    config.Addr,
		Handler: router,
	}
}

// Viewer views the swagger file
type Viewer struct {
	Path string
}

// Mount mounts the editor
func (e *Viewer) Mount(r chi.Router) {
	r.Get("/*", e.serve)
	r.Get("/swagger.spec", e.load)
}

func (e *Viewer) serve(w http.ResponseWriter, r *http.Request) {
	handler := http.FileServer(parcello.ManagerAt("viewer"))
	handler.ServeHTTP(w, r)
}

func (e *Viewer) load(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, e.Path)
}
