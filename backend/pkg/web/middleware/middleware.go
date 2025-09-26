package middleware

import (
	"encoding/json"
	"net/http"
)

func MaxBodyLength(length int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rr := new(http.Request)
			*rr = *r
			rr.Body = http.MaxBytesReader(w, rr.Body, length)
			next.ServeHTTP(w, rr)
		})
	}
}

func RenderJSON(w http.ResponseWriter, resp any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	raw, _ := json.Marshal(resp)
	w.Write(raw)
}
