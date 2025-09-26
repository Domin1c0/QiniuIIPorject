package web

import (
	"net/http"

	storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"
	"github.com/LTSlw/QiniuIIPorject/backend/pkg/web/middleware"
	"github.com/go-chi/chi/v5"
)

func (s *Server) basicRouters() *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", s.health)
	r.Post("/register", s.register)
	r.Post("/login", s.login)
	r.With(middleware.AuthRequired(s.db)).Put("/logout", s.logout)
	return r
}

func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	// TODO build req/resp model
	username := r.FormValue("username")
	password := r.FormValue("password")

	// TODO IMPORTANT hash password
	if username == "" || password == "" {
		middleware.Error(w, http.StatusBadRequest, ErrRequestMissingFields)
	}
	user, err := s.db.AddUser(&storage.User{Name: username, Password: password})
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, err)
	}
	middleware.RenderJSON(w, user)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	if username == "" || password == "" {
		middleware.Error(w, http.StatusBadRequest, ErrRequestMissingFields)
		return
	}

	user, err := s.db.GetUserByName(username)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, err)
		return
	}
	if user == nil || user.Password != password {
		middleware.Error(w, http.StatusUnauthorized, ErrUserOrPasswordWrong)
		return
	}

	session, err := s.db.UserLogin(user.Id, password)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, err)
		return
	}
	middleware.RenderJSON(w, session)
}

func (s *Server) logout(w http.ResponseWriter, r *http.Request) {
	session := middleware.GetUserSession(r)
	s.db.UserLogout(session)
}
