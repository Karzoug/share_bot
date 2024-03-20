package middleware

import (
	"net/http"
)

// HeaderAuth implements a simple middleware handler for adding http header auth to a route.
func HeaderAuth(headerKey string, authFn func(string) bool) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if v := r.Header.Get(headerKey); !authFn(v) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
