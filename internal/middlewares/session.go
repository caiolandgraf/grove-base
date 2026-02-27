package middlewares

import (
	"net/http"

	"github.com/alexedwards/scs/v2"
)

// SessionMiddleware integrates SCS with HTTP handlers
func SessionMiddleware(
	sessionManager *scs.SessionManager,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return sessionManager.LoadAndSave(next)
	}
}

// AuthRequired middleware - checks if user is logged in
func AuthRequired(
	sessionManager *scs.SessionManager,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sessionManager.GetString(r.Context(), "user_id")

			if userID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// User authenticated, continue
			next.ServeHTTP(w, r)
		})
	}
}

// GuestOnly middleware - only for unauthenticated users
func GuestOnly(
	sessionManager *scs.SessionManager,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := sessionManager.GetString(r.Context(), "user_id")

			if userID != "" {
				http.Error(w, "Already authenticated", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
