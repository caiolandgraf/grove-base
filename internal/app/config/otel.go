package config

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
)

// OtelShutdown is a function to call on application shutdown.
type OtelShutdown func(ctx context.Context) error

// InitOtel initializes OpenTelemetry with an OTLP HTTP exporter.
func InitOtel(ctx context.Context) (OtelShutdown, error) {
	if !Env.OtelEnabled {
		slog.Info("OpenTelemetry tracing disabled")
		return func(_ context.Context) error { return nil }, nil
	}

	serviceName := Env.OtelServiceName
	otelEndpoint := Env.OtelOTLPEndpoint

	sampleRatio := Env.OtelTraceSampleRatio
	if sampleRatio < 0 {
		sampleRatio = 0
	}
	if sampleRatio > 1 {
		sampleRatio = 1
	}

	otelEndpoint = strings.TrimPrefix(otelEndpoint, "http://")
	otelEndpoint = strings.TrimPrefix(otelEndpoint, "https://")

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion("2.0.0"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel resource: %w", err)
	}

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(otelEndpoint),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTel exporter: %w", err)
	}

	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(sampleRatio))

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sampler),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	slog.Info("OpenTelemetry initialized",
		"service", serviceName,
		"endpoint", otelEndpoint,
		"sample_ratio", sampleRatio,
	)

	shutdown := func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		return tp.Shutdown(ctx)
	}

	return shutdown, nil
}
