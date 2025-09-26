package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"
	"github.com/LTSlw/QiniuIIPorject/backend/pkg/web/middleware"
	"github.com/LTSlw/QiniuIIPorject/backend/pkg/web/model"
	"github.com/go-chi/chi/v5"
)

func (s *Server) chatRouters() *chi.Mux {
	r := chi.NewRouter()

	r.With(middleware.AuthRequired(s.db)).Get("/sessions", s.getSessions)
	r.With(middleware.AuthRequired(s.db)).Get("/{sessionID}/messages", s.getSessionMessages)

	r.With(middleware.AuthRequired(s.db)).Post("/new", s.postNewMessages)
	r.With(middleware.AuthRequired(s.db)).Post("/{sessionID}/send", s.postMessages)
	return r
}

func (s *Server) getSessions(w http.ResponseWriter, r *http.Request) {
	userSession := middleware.GetUserSession(r)

	urlNum := r.URL.Query().Get("num")
	num := 16
	if urlNum != "" {
		if n, err := strconv.Atoi(urlNum); err == nil {
			num = n
		} else {
			middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
			return
		}
	}
	urlFrom := r.URL.Query().Get("from")
	from := time.Now()
	if urlFrom != "" {
		if ts, err := strconv.ParseInt(urlFrom, 10, 64); err == nil {
			from = time.Unix(ts, 0)
		} else {
			middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
			return
		}
	}

	sessions, err := s.db.GetSessionsByUserID(userSession.UserId, num, from)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	middleware.RenderJSON(w, sessions)
}

func (s *Server) getSessionMessages(w http.ResponseWriter, r *http.Request) {
	userSession := middleware.GetUserSession(r)
	urlSessionId, err := strconv.Atoi(chi.URLParam(r, "sessionID"))
	if err != nil {
		middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
		return
	}

	// Check session ownership
	session, err := s.db.GetSessionByID(urlSessionId)
	if err != nil {
		if err == storage.ErrNotFound {
			middleware.Error(w, http.StatusNotFound, ErrNotFound)
		} else {
			middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		}
		return
	}
	if session.UserId != userSession.UserId {
		middleware.Error(w, http.StatusForbidden, ErrPermissionDenied)
		return
	}

	messages, err := s.db.GetMessagesBySessionID(urlSessionId)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	middleware.RenderJSON(w, messages)
}

func (s *Server) postNewMessages(w http.ResponseWriter, r *http.Request) {
	// TODO
	userSession := middleware.GetUserSession(r)

	var req model.RequestNewChat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
		return
	}

	session, err := s.db.AddSession(storage.Session{
		UserId: userSession.UserId,
	})
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	message, err := s.db.AddMessage(storage.Message{
		SessionId: session.Id,
		Role:      "user",
		Content:   req.Message.Content,
	})
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	// Should build a response model
	resp := storage.SessionWithMessages{
		Session:  session,
		Messages: []storage.Message{message},
	}

	middleware.RenderJSON(w, resp)
}

func (s *Server) postMessages(w http.ResponseWriter, r *http.Request) {
	// TODO
	userSession := middleware.GetUserSession(r)
	urlSessionId, err := strconv.Atoi(chi.URLParam(r, "sessionID"))
	if err != nil {
		middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
		return
	}
	var req model.RequestSend
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
		return
	}
	// Check session ownership
	session, err := s.db.GetSessionByID(urlSessionId)
	if err != nil {
		if err == storage.ErrNotFound {
			middleware.Error(w, http.StatusNotFound, ErrNotFound)
		} else {
			middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		}
		return
	}
	if session.UserId != userSession.UserId {
		middleware.Error(w, http.StatusForbidden, ErrPermissionDenied)
		return
	}

	// TODO IMPORTANT call LLM, now only add a new message to message list
	message, err := s.db.AddMessage(storage.Message{
		SessionId: session.Id,
		Role:      "user",
		Content:   req.Message.Content,
	})
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	middleware.RenderJSON(w, message)
}
