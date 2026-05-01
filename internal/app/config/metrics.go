package config

import (
	"log/slog"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprometheus "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"go.opentelemetry.io/otel"
)

// InitMetrics initializes the OpenTelemetry MeterProvider with a Prometheus exporter.
// This bridges OTel metrics (from otelhttp, otelgorm, etc.) to Prometheus format.
// Returns the HTTP handler to be registered at /metrics.
func InitMetrics() (http.Handler, error) {
	exporter, err := otelprometheus.New()
	if err != nil {
		return nil, err
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(exporter),
	)

	otel.SetMeterProvider(mp)

	slog.Info("Prometheus metrics initialized", "endpoint", "/metrics")

	return promhttp.Handler(), nil
}
