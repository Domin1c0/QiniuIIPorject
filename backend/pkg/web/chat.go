package web

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) chatRouters() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/sessions", s.getSessions)
	return r
}

func (s *Server) getSessions(w http.ResponseWriter, r *http.Request) {

}
