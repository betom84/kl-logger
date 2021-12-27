package metrics

import (
	"fmt"
	"net/http"
	"time"

	"github.com/betom84/kl-logger/klimalogg/frames"
	"github.com/betom84/kl-logger/repository"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	registry = prometheus.DefaultRegisterer

	consoleFrames = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "klimalogg_console_frames",
			Help: "KlimaLogg console frames received via usb tranceiver (GetFrame)",
		},
		[]string{"loggerID", "deviceID", "typeID"},
	)

	temperature = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "klimalogg_sensor_temperature",
			Help: "Latest temperature values received from KlimaLogg console by sensor",
		},
		[]string{"sensor"},
	)

	humidity = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "klimalogg_sensor_humidity",
			Help: "Latest humidity values received from KlimaLogg console by sensor",
		},
		[]string{"sensor"},
	)

	signalQuality = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "klimalogg_console_signal_quality",
			Help: "Current signal quality between KlimaLogg console and usb tranceiver",
		},
	)
)

func init() {
	registry.MustRegister(consoleFrames, temperature, humidity, signalQuality)
}

func HttpMiddleware(next http.Handler) http.Handler {

	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of http requests by method, endpoint and statusCode",
		},
		[]string{"method", "endpoint", "statusCode"},
	)

	requestDuration := prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name:       "http_request_duration_milliseconds",
			Help:       "Http request duration in milliseconds by method, endpoint and statusCode",
			Objectives: map[float64]float64{0.9: 0.01, 0.95: 0.005, 0.99: 0.001},
		},
		[]string{"method", "endpoint", "statusCode"},
	)

	registry.MustRegister(requestsTotal, requestDuration)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		wrapped := middleware.NewWrapResponseWriter(w, 0)

		defer func() {
			labels := prometheus.Labels{
				"method":     r.Method,
				"endpoint":   r.URL.Path,
				"statusCode": fmt.Sprint(wrapped.Status()),
			}

			requestsTotal.With(labels).Inc()
			requestDuration.With(labels).Observe(float64(time.Now().UnixMilli() - startTime.UnixMilli()))
		}()

		next.ServeHTTP(wrapped, r)
	})
}

func HttpHandler() http.Handler {

	return promhttp.Handler()
}

// todo, implement repository as prometheus collector
func KlimaloggCurrentValuesPublisher() chan<- interface{} {
	c := make(chan interface{})

	go func() {
		for {
			i := <-c

			switch t := i.(type) {

			case repository.WeatherSample:
				for i := t.SensorMin(); i <= t.SensorMax(); i++ {
					if !t.IsSensorActive(i) {
						continue
					}

					temperature.With(prometheus.Labels{"sensor": fmt.Sprint(i)}).Set(float64(t.Temperature(i)))
					humidity.With(prometheus.Labels{"sensor": fmt.Sprint(i)}).Set(float64(t.Humidity(i)))
				}
				signalQuality.Set(float64(t.SignalQuality()))
			}
		}
	}()

	return c
}

func CollectConsoleFrame(frame frames.GetFrame) {
	consoleFrames.With(prometheus.Labels{
		"loggerID": fmt.Sprint(frame.LoggerID()),
		"deviceID": fmt.Sprintf("%04x", frame.DeviceID()),
		"typeID":   frame.TypeID().String(),
	}).Inc()
}
