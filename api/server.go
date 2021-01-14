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
	router     chi.Router
	repository Repository
}

// NewServer to serve api endpoints
func NewServer(repository Repository) *Server {
	s := &Server{
		repository: repository,
		router:     chi.NewRouter(),
	}

	s.router.Mount("/debug", middleware.Profiler())
	s.router.Get("/current", getCurrentWeather(s.repository))

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
