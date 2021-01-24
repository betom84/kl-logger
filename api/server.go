package api

import (
	"net/http"

	"github.com/betom84/kl-logger/repository"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Server for klimalogg api
type Server struct {
	router chi.Router
}

// NewServer to serve api endpoints
func NewServer(repo repository.Repository, transceiver Traceable) *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	s.router.Mount("/debug/pprof", middleware.Profiler())

	if transceiver != nil {
		s.router.Get("/debug/transceiver/trace", GetTransceiverTrace(transceiver))
	}

	s.router.Get("/weather", GetWeather(repo))
	s.router.Get("/weather/{sensor:[0-8]}", GetWeatherBySensor(repo))

	s.router.Get("/config", GetConfig(repo))
	s.router.Get("/config/{sensor:[0-8]}", GetConfigBySensor(repo))

	s.router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		w.Write([]byte("not found"))
	})

	return s
}

// ServeHTTP requests
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}
