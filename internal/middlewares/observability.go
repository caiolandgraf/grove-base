package middlewares

import (
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.39.0"
	"go.opentelemetry.io/otel/trace"
)

// RouteTagMiddleware reads the matched route pattern from Go 1.22+ (r.Pattern)
// and sets it as the http.route attribute on both the otelhttp labeler (so
// Prometheus metrics include the route label) and the current trace span (so
// Jaeger shows the route in the span name and attributes).
//
// This middleware MUST be registered AFTER otelhttp.NewMiddleware so that the
// labeler context is available.
func RouteTagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Serve the request first so the downstream mux has a chance to
		// populate r.Pattern (fuego uses Go's standard ServeMux under the hood).
		// However, fuego registers each route individually on the mux, so the
		// pattern is already set BEFORE this middleware runs. We can read it
		// immediately.
		route := normalizePattern(r.Pattern)
		if route == "" {
			// Fallback: use the raw path (better than nothing).
			route = r.URL.Path
		}

		// --- Set on otelhttp Labeler (affects Prometheus metrics) ---
		if labeler, ok := otelhttp.LabelerFromContext(r.Context()); ok {
			labeler.Add(semconv.HTTPRouteKey.String(route))
		}

		// --- Set on the current trace span (affects Jaeger) ---
		span := trace.SpanFromContext(r.Context())
		if span.IsRecording() {
			span.SetAttributes(
				attribute.String(string(semconv.HTTPRouteKey), route),
			)
			// Rename span to "METHOD /route" for readability in Jaeger.
			span.SetName(r.Method + " " + route)
		}

		next.ServeHTTP(w, r)
	})
}

// normalizePattern cleans up the Go 1.22+ ServeMux pattern.
// Patterns may contain a method prefix like "GET /api/v1/users/{id}"
// or "GET example.com/path". We strip the method and host, keeping only the
// path portion so that the route label stays clean and low-cardinality.
func normalizePattern(pattern string) string {
	if pattern == "" {
		return ""
	}

	// Remove leading/trailing whitespace.
	pattern = strings.TrimSpace(pattern)

	// Go 1.22+ patterns may start with "METHOD " (e.g. "GET /foo").
	// Strip the method prefix if present.
	if idx := strings.Index(pattern, " "); idx != -1 {
		pattern = strings.TrimSpace(pattern[idx+1:])
	}

	// Strip optional host prefix (e.g. "example.com/path" → "/path").
	if !strings.HasPrefix(pattern, "/") {
		if idx := strings.Index(pattern, "/"); idx != -1 {
			pattern = pattern[idx:]
		}
	}

	// Remove trailing "{$}" exact-match marker that Go's mux may include.
	pattern = strings.TrimSuffix(pattern, "{$}")
	pattern = strings.TrimRight(pattern, "/")

	if pattern == "" {
		pattern = "/"
	}

	return pattern
}
