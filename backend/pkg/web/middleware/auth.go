package middleware

import (
	"context"
	"net/http"
	"strings"

	storage "github.com/LTSlw/QiniuIIPorject/backend/pkg/storge"
)

type contextKey string

const userSessionKey contextKey = "userSession"

func Auth(db *storage.Storage, required bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				if required {
					Error(w, http.StatusUnauthorized, ErrMissingToken)
					return
				} else {
					next.ServeHTTP(w, r)
					return
				}
			}

			const prefix = "Bearer "
			authHeader = strings.TrimPrefix(authHeader, prefix)

			session, err := db.GetUserSessionByID(authHeader)
			if err != nil {
				Error(w, http.StatusInternalServerError, err)
				return
			}
			if session == nil {
				Error(w, http.StatusUnauthorized, ErrInvalidToken)
				return
			}

			ctx := context.WithValue(r.Context(), userSessionKey, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AuthRequired(db *storage.Storage) func(http.Handler) http.Handler {
	return Auth(db, true)
}

func AuthOptional(db *storage.Storage) func(http.Handler) http.Handler {
	return Auth(db, false)
}

func GetUserSession(r *http.Request) *storage.UserSession {
	val := r.Context().Value(userSessionKey)
	if session, ok := val.(*storage.UserSession); ok {
		return session
	}
	return nil
}
