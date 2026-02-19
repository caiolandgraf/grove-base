package routes

import (
	"net/http"

	"github.com/caiolandgraf/go-project-base/internal/container"
	"github.com/caiolandgraf/go-project-base/internal/middlewares"
	"github.com/go-fuego/fuego"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// SetupRoutes configures all routes using the container
func SetupRoutes(
	s *fuego.Server,
	c *container.Container,
	metricsHandler http.Handler,
) {
	// OpenTelemetry Middleware
	fuego.Use(s, otelhttp.NewMiddleware("go-project-base"))

	// Route tag middleware — reads r.Pattern (Go 1.22+) and sets http.route
	// on the otelhttp labeler (Prometheus metrics) and span (Jaeger traces).
	fuego.Use(s, middlewares.RouteTagMiddleware)

	// CORS Middleware global
	fuego.Use(s, middlewares.CORSMiddleware(middlewares.DefaultCORSConfig()))

	// Session Middleware global
	fuego.Use(s, middlewares.SessionMiddleware(c.SessionManager))

	// Health check
	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed(c))

	// Prometheus metrics endpoint
	fuego.GetStd(s, "/metrics", metricsHandler.ServeHTTP)

	// API v1
	api := fuego.Group(s, "/api/v1")

	// ========== AUTH ROUTES (Public) ==========
	auth := fuego.Group(api, "/auth")
	fuego.Post(auth, "/login", c.AuthController.Login)
	fuego.Post(auth, "/register", c.UserController.CreateUser)
	fuego.Post(auth, "/logout", c.AuthController.Logout)
	fuego.Get(auth, "/me", c.AuthController.Me)

	// ========== USER ROUTES (Protected) ==========
	users := fuego.Group(api, "/users")
	// Auth Middleware if needed for all user routes
	// users.Use(middlewares.AuthRequired(c.SessionManager))
	// GET /users
	fuego.Get(users, "/", c.UserController.ListUsers)
	// POST /users
	fuego.Post(users, "/", c.UserController.CreateUser)
	// GET /users/{user_id}
	fuego.Get(users, "/{user_id}", c.UserController.GetUser)
	// PUT /users/{user_id}
	fuego.Put(users, "/{user_id}", c.UserController.UpdateUser)
	// DELETE /users/{user_id}
	fuego.Delete(users, "/{user_id}", c.UserController.DeleteUser)
}
