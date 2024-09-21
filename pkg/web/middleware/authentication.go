package middleware

import (
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/web/session"
)

func LoginRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.FromContext(r.Context())
		if sess == nil || !sess.LoggedIn {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}
