package api

import (
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Repository to retrieve weather data
type Repository interface {
	Temperature() float32
	Humidity() uint
	LastUpdate() time.Time
}

// Server for klimalogg api
type Server struct {
	router chi.Router
}

// NewServer to serve api endpoints
func NewServer() *Server {
	s := &Server{
		router: chi.NewRouter(),
	}

	s.router.Mount("/debug", middleware.Profiler())
	s.router.Get("/weather", getCurrentWeather())
	s.router.Get("/config", getCurrentConfig())

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
