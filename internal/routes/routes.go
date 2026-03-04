package routes

import (
	"github.com/caiolandgraf/grove-base/internal/app"
	"github.com/caiolandgraf/grove-base/internal/controllers"
	"github.com/caiolandgraf/grove-base/internal/middleware"
	"github.com/go-fuego/fuego"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// ──────────────────────────────────────────────
// Health check DTOs
// ──────────────────────────────────────────────

type HealthCheckResponse struct {
	Status string `json:"status"`
}

type HealthCheckDetailedResponse struct {
	Status   string            `json:"status"`
	Services map[string]string `json:"services"`
}

// ──────────────────────────────────────────────
// Routes
// ──────────────────────────────────────────────

// SetupRoutes configures all routes using app globals
func SetupRoutes(s *fuego.Server) {
	authController := controllers.NewAuthController(app.Session)

	// OpenTelemetry Middleware
	fuego.Use(s, otelhttp.NewMiddleware("grove-app"))

	// Route tag middleware — reads r.Pattern (Go 1.22+) and sets http.route
	// on the otelhttp labeler (Prometheus metrics) and span (Jaeger traces).
	fuego.Use(s, middleware.RouteTagMiddleware)

	// CORS Middleware global
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))

	// Session Middleware global
	fuego.Use(s, middleware.SessionMiddleware(app.Session))

	// Health check
	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed)

	// Prometheus metrics endpoint
	fuego.GetStd(s, "/metrics", app.Metrics.ServeHTTP)

	// API v1
	api := fuego.Group(s, "/api/v1")

	// ========== AUTH ROUTES (Public) ==========
	auth := fuego.Group(api, "/auth")
	fuego.Post(auth, "/login", authController.Login)
	fuego.Post(auth, "/register", authController.Register)
	fuego.Post(auth, "/logout", authController.Logout)
	fuego.Get(auth, "/me", authController.Me)

	// ========== USER ROUTES ==========
	users := fuego.Group(api, "/users")
	// Auth Middleware if needed for all user routes
	// fuego.Use(users, middleware.AuthRequired(app.Session))
	fuego.Get(users, "/", controllers.ListUsers)
	fuego.Post(users, "/", controllers.CreateUser)
	fuego.Get(users, "/{user_id}", controllers.GetUser)
	fuego.Put(users, "/{user_id}", controllers.UpdateUser)
	fuego.Delete(users, "/{user_id}", controllers.DeleteUser)
}

// ──────────────────────────────────────────────
// Health check handlers
// ──────────────────────────────────────────────

func healthCheck(c fuego.ContextNoBody) (*HealthCheckResponse, error) {
	return &HealthCheckResponse{
		Status: "Ok",
	}, nil
}

func healthCheckDetailed(
	c fuego.ContextNoBody,
) (*HealthCheckDetailedResponse, error) {
	return &HealthCheckDetailedResponse{
		Status:   "OK",
		Services: app.HealthCheck(),
	}, nil
}
