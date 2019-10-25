package service

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Config is the service config
// stride:generate config
type Config struct {
	// stride:generate addr
	Addr string
}

// Route represents a mountable route
// stride:generate route
type Route interface {
	Mount(r chi.Router)
}

// NewServer creates a new server
// stride:generate new-server
func NewServer(config *Config) *http.Server {
	router := chi.NewRouter()
	router.Use(middleware.StripSlashes)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.NoCache)
	router.Use(middleware.Logger)

	routes := []Route{
	  {{- range .controllers }}
    // stride:generate
	  &{{ .Name | camelize }}API{},
	  {{- end }}
  }

	for _, route := range routes {
		route.Mount(router)
	}

	server := &http.Server{
		Addr:    config.Addr,
		Handler: router,
	}

	return server
}
