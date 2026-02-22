package middleware

import (
	"net/http"

	"google-sso/auth"
)

// RequireAuth wraps a handler and redirects unauthenticated requests to /login.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !auth.IsAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		next.ServeHTTP(w, r)
	})
}
