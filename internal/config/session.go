package config

import (
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func InitSessionManager(redisPool *redis.Pool) *scs.SessionManager {
	sessionManager := scs.New()

	// Configure Redis as store
	sessionManager.Store = redisstore.New(redisPool)

	// Security settings
	sessionManager.Lifetime = 24 * time.Hour      // Session lasts 24h
	sessionManager.IdleTimeout = 20 * time.Minute // Inactivity timeout
	sessionManager.Cookie.Name = "session_id"
	sessionManager.Cookie.HttpOnly = true // Not accessible via JavaScript
	sessionManager.Cookie.Persist = true  // Persist after closing browser
	sessionManager.Cookie.SameSite = 1    // Strict (CSRF protection)
	sessionManager.Cookie.Secure = false  // true in production (HTTPS)
	sessionManager.Cookie.Path = "/"

	return sessionManager
}
