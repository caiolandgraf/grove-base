package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// OtelShutdown is a function to call on application shutdown
type OtelShutdown func(ctx context.Context) error

// InitOtel initializes OpenTelemetry with an OTLP HTTP exporter.
// Returns a shutdown function that should be called on application exit.
func InitOtel(ctx context.Context) (OtelShutdown, error) {
	serviceName := os.Getenv("OTEL_SERVICE_NAME")
	if serviceName == "" {
		serviceName = "go-project-base"
	}

	otelEndpoint := os.Getenv("OTEL_EXPLOERER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "localhost:4318"
	}

	// WithEndpoint expects host:port, not a full URL — strip scheme if present
	otelEndpoint = strings.TrimPrefix(otelEndpoint, "http://")
	otelEndpoint = strings.TrimPrefix(otelEndpoint, "https://")

	// Resource with service info
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("1.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	// Exporter OTLP via HTTP
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otelEndpoint),
		otlptracehttp.WithInsecure(), // Remove in production with TLS
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel exporter: %w", err)
	}

	// Tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Register as global TracerProvider so otelhttp and otelgorm pick it up
	otel.SetTracerProvider(tp)

	// Set global propagator for trace context propagation across services
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	slog.Info("OpenTelemetry initialized",
		"service", serviceName,
		"endpoint", otelEndpoint,
	)

	// Return shutdown function
	shutdown := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}

	return shutdown, nil
}
