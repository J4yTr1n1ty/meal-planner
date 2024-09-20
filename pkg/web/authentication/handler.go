package authentication

import "net/http"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, post-check=0, pre-check=0")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	http.ServeFile(w, r, "static/login.html")
}

func (h *Handler) LoginPost(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement login
	w.Header().Set("HX-Redirect", "/")
	w.WriteHeader(http.StatusFound)
}
