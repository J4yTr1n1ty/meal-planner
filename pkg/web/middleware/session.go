package middleware

import (
	"context"
	"log"
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/web/session"
)

func Session(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("session middleware")
		sess := session.LoadOrNew(r)
		ctx := context.WithValue(r.Context(), session.ContextKey, sess)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
