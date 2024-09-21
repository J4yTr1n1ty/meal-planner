package authentication

import (
	"net/http"

	"github.com/J4yTr1n1ty/meal-planner/pkg/config"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/htmx"
	"github.com/J4yTr1n1ty/meal-planner/pkg/web/session"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sess := session.FromContext(r.Context())
		if sess != nil && sess.LoggedIn {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		http.ServeFile(w, r, "static/login.html")
	}
}

func (h *Handler) LoginPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			htmx.RenderError(w, http.StatusOK, err.Error())
		}
		password := r.Form.Get("password")
		if password == config.Password {
			session.LoginUser(w, r)
			w.Header().Set("HX-Redirect", "/")
			w.WriteHeader(http.StatusOK)
			return
		} else {
			htmx.RenderError(w, http.StatusOK, "Wrong password")
		}
	}
}
