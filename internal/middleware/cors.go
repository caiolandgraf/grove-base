package middleware

import (
	"net/http"
	"strings"

	"github.com/caiolandgraf/grove-base/internal/config"
)

// CORSConfig contains the CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	AllowCredentials bool
	MaxAge           string
}

// DefaultCORSConfig returns the default configuration (adjust for production!)
func DefaultCORSConfig() CORSConfig {
	return CORSConfig{
		AllowedOrigins: strings.Split(config.Env.CORSAllowedOrigins, ","),
		AllowedMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowedHeaders: []string{
			"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin",
		},
		AllowCredentials: true,
		MaxAge:           "86400",
	}
}

// setCORSHeaders writes the CORS headers to the response
func setCORSHeaders(
	w http.ResponseWriter,
	r *http.Request,
	cfg CORSConfig,
	allowedOriginsSet map[string]bool,
) {
	origin := r.Header.Get("Origin")

	if allowedOriginsSet[origin] {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}

	if cfg.AllowCredentials {
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.Header().
		Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
	w.Header().
		Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
	w.Header().Set("Access-Control-Max-Age", cfg.MaxAge)
}

// buildOriginsSet creates the set of allowed origins
func buildOriginsSet(cfg CORSConfig) map[string]bool {
	set := make(map[string]bool)
	for _, o := range cfg.AllowedOrigins {
		set[strings.TrimSpace(o)] = true
	}
	return set
}

// CORSMiddleware returns an HTTP middleware for CORS (headers in actual responses)
func CORSMiddleware(cfg CORSConfig) func(http.Handler) http.Handler {
	allowedOriginsSet := buildOriginsSet(cfg)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			setCORSHeaders(w, r, cfg, allowedOriginsSet)

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// CORSPreflightHandler returns a handler to respond to preflight OPTIONS requests
func CORSPreflightHandler(cfg CORSConfig) http.HandlerFunc {
	allowedOriginsSet := buildOriginsSet(cfg)

	return func(w http.ResponseWriter, r *http.Request) {
		setCORSHeaders(w, r, cfg, allowedOriginsSet)
		w.WriteHeader(http.StatusNoContent)
	}
}
