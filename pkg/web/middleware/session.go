package middleware

import (
	"net/http"
)

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Add Session loading logic
		next.ServeHTTP(w, r)
	})
}
