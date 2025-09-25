package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) basicRouters() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", s.health)
	return r
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
