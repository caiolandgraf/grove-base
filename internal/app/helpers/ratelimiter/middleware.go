package ratelimiter

import (
	"net/http"

	"github.com/go-fuego/fuego"
)

// Middleware returns HTTP middleware that applies the same limits as Check.
// Use Check inside handlers when you need per-route or per-operation control.
func (l *Limiter) Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if err := l.Check(w, r); err != nil {
				fuego.SendError(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
