package api

import (
	"encoding/json"
	"net/http"
	"time"
)

func getCurrentWeather(repository Repository) handlerFuncWithError {
	type response struct {
		UpdatedAt   string  `json:"updatedAt"`
		Temperature float32 `json:"temperature"`
		Humidity    uint    `json:"humidity"`
	}

	return func(w http.ResponseWriter, r *http.Request) error {
		response := response{
			Temperature: repository.Temperature(),
			Humidity:    repository.Humidity(),
			UpdatedAt:   repository.LastUpdate().Format(time.RFC3339),
		}

		json.NewEncoder(w).Encode(response)
		w.Header().Set("Content-Type", "application/json")

		return nil
	}
}
