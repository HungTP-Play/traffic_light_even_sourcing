package metrics

import (
	"bytes"

	"github.com/gofiber/fiber/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/common/expfmt"
)

var (
	// This will use in middleware output format => requests_total{route="/",path="/",status="200",method="GET"}
	RequestsTotal = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "traffic_light_projection_requests_total",
			Help: "Number of requests received",
		},
		[]string{"method", "path"},
	)

	TrafficLightState = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "traffic_light_projection_light_state",
			Help: "The current state of a traffic light (1=RED, 2=GREEN, 3=YELLOW)",
		},
		[]string{"id", "location_lat", "location_long"},
	)
)

func PrometheusMiddleware(c *fiber.Ctx) error {
	RequestsTotal.WithLabelValues(c.Method(), c.Path()).Observe(1)
	return c.Next()
}

func getPrometheusMetrics() (string, error) {
	var metrics bytes.Buffer
	reg := prometheus.DefaultRegisterer.(*prometheus.Registry)
	metricFamilies, err := reg.Gather()
	if err != nil {
		return "", err
	}

	encoder := expfmt.NewEncoder(&metrics, expfmt.FmtText)
	for _, mf := range metricFamilies {
		if err := encoder.Encode(mf); err != nil {
			return "", err
		}
	}

	return metrics.String(), nil
}

func Metrics(c *fiber.Ctx) error {
	metrics, err := getPrometheusMetrics()
	if err != nil {
		return c.Status(500).SendString("Failed to collect metrics")
	}
	return c.Type("text/plain").SendString(metrics)
}

func SetTrafficLightState(id, locationLat, locationLong string, color int) {
	TrafficLightState.WithLabelValues(id, locationLat, locationLong).Set(float64(color))
}
