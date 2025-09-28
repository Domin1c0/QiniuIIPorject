package web

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/LTSlw/QiniuIIPorject/backend/pkg/llm"
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

// 	middleware.RenderJSON(w, llmMessage)
// }

// public handle function chat message
// used by both new session and existing session
// Modified from commented postNewMessages and postMessages functions
func (s *Server) handleChatMessage(
	w http.ResponseWriter,
	r *http.Request,
	userID int,
	sessionID *int, // == nil for new session, != nil for existing session
	content string,
) {
	var session storage.Session
	var err error

	// if sessionID is nil, create a new session
	// otherwise get the session and check ownership
	if sessionID == nil {
		session, err = s.db.AddSession(storage.Session{UserId: userID})
		if err != nil {
			middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
			return
		}
	} else {
		session, err = s.db.GetSessionByID(*sessionID)
		if err != nil {
			if err == storage.ErrNotFound {
				middleware.Error(w, http.StatusNotFound, ErrNotFound)
			} else {
				middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
			}
			return
		}
		if session.UserId != userID {
			middleware.Error(w, http.StatusForbidden, ErrPermissionDenied)
			return
		}
	}

	// TODO Insert a system message(Contain character prompt) to the new session

	// Insert user message to db
	userMessage, err := s.db.AddMessage(storage.Message{
		SessionId: session.Id,
		Role:      "user",
		Content:   content,
	})
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}

	// Select context messages based on token limit
	messages, err := s.db.GetMessagesBySessionID(session.Id)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	messages = append(messages, userMessage)

	contextMessages, err := llm.SelectMessage(
		storage.SessionWithMessages{
			Session:  session,
			Messages: messages,
		},
		s.llmModel.MaxToken,
		s.llmModel,
	)
	if err != nil {
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	s.logger.Debug().Any("contextMessages", contextMessages).Msg("Selected context messages")

	// Call LLM
	llmMessage, err := llm.CallLLM(
		s.llmModel,
		contextMessages,
	)
	if err != nil {
		s.logger.Error().Err(err).Msg("CallLLM error")
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}
	llmMessage.SessionId = session.Id

	// Insert llm message to db
	llmMessage, err = s.db.AddMessage(llmMessage)
	if err != nil {
		s.logger.Error().Err(err).Msg("Failed to insert llm message to db")
		middleware.Error(w, http.StatusInternalServerError, ErrInternalError)
		return
	}

	// only response the latest message by llm
	// middleware.RenderJSON(w, llmMessage)
	var resp model.ChatResponse
	resp.SessionId = session.Id
	resp.Message = llmMessage
	middleware.RenderJSON(w, resp)
}

// Post a message to new session
func (s *Server) postNewMessages(w http.ResponseWriter, r *http.Request) {
	// TODO IMPORTANT system prompt message (contain character prompt) in new session
	userSession := middleware.GetUserSession(r)

	var req model.RequestNewChat
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		middleware.Error(w, http.StatusBadRequest, ErrRequestInvalid)
		return
	}

	s.handleChatMessage(w, r, userSession.UserId, nil, req.Message.Content)
}

// Post a message to existing session
func (s *Server) postMessages(w http.ResponseWriter, r *http.Request) {
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

	s.handleChatMessage(w, r, userSession.UserId, &urlSessionId, req.Message.Content)
}
