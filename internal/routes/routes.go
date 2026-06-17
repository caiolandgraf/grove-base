package routes

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/grove-base/internal/app/config"
	"github.com/caiolandgraf/grove-base/internal/app/middleware"
	"github.com/caiolandgraf/grove-base/internal/modules"
	"github.com/go-fuego/fuego"
	"github.com/gomodule/redigo/redis"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"gorm.io/gorm"
)

func SetupRoutes(
	s *fuego.Server,
	db *gorm.DB,
	redisPool *redis.Pool,
	session *scs.SessionManager,
	metricsHandler http.Handler,
) {
	if config.Env.OtelEnabled || config.Env.MetricsEnabled {
		fuego.Use(s, otelhttp.NewMiddleware(config.Env.OtelServiceName))
		fuego.Use(s, middleware.RouteTagMiddleware)
	}

	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))
	fuego.Use(s, middleware.SessionMiddleware(session))

	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed(db, redisPool))

	if metricsHandler != nil {
		fuego.GetStd(s, "/metrics", metricsHandler.ServeHTTP)
	}

	api := fuego.Group(s, "/api/v1")
	modules.Mount(api, modules.Boot{
		DB:        db,
		Session:   session,
		RateLimit: config.RateLimitSettings(),
	})
}
