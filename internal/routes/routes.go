package routes

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/caiolandgraf/go-project-base/internal/modules"
	"github.com/caiolandgraf/go-project-base/internal/app/middleware"
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
	fuego.Use(s, otelhttp.NewMiddleware("go-project-base"))
	fuego.Use(s, middleware.RouteTagMiddleware)
	fuego.Use(s, middleware.CORSMiddleware(middleware.DefaultCORSConfig()))
	fuego.Use(s, middleware.SessionMiddleware(session))

	fuego.Get(s, "/", healthCheck)
	fuego.Get(s, "/health", healthCheckDetailed(db, redisPool))
	fuego.GetStd(s, "/metrics", metricsHandler.ServeHTTP)

	api := fuego.Group(s, "/api/v1")
	modules.Mount(api, modules.Boot{DB: db, Session: session})
}
