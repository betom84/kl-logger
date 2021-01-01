package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// Repository to retrieve weather data
type Repository interface {
	Temperature() float32
	Humidity() uint
	LastUpdate() time.Time
}

// Server for klimalogg api
type Server struct {
	router     *router
	repository Repository
}

// NewServer to serve api endpoints
func NewServer(repository Repository) *Server {
	s := &Server{
		router:     newRouter(),
		repository: repository,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.router.get("/current", getCurrentWeather(s.repository))

	s.router.fallback = func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(404)
		w.Write([]byte("not found"))

		return nil
	}
}

// ServeHTTP requests
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.router.serveHTTP(w, r)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":   err,
			"request": r,
		}).Error("error handling http request")

		w.Header().Add("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{ "error": "%s" }`, err.Error())))
	}
}
