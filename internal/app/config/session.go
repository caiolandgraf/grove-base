package config

import (
	"time"

	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gomodule/redigo/redis"
)

func InitSessionManager(redisPool *redis.Pool) *scs.SessionManager {
	sessionManager := scs.New()

	sessionManager.Store = redisstore.New(redisPool)

	sessionManager.Lifetime = 24 * time.Hour
	sessionManager.IdleTimeout = 20 * time.Minute
	sessionManager.Cookie.Name = "session_id"
	sessionManager.Cookie.HttpOnly = true
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.SameSite = 1
	sessionManager.Cookie.Secure = false
	sessionManager.Cookie.Path = "/"

	if Env.Environment == "production" {
		sessionManager.Cookie.Secure = true
		sessionManager.Cookie.Domain = Env.CookieDomain
	}

	return sessionManager
}
