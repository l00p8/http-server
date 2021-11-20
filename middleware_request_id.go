package xserver

import (
	"net/http"

	chiMiddleware "github.com/go-chi/chi/middleware"
)

func xRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		reqID := r.Context().Value(chiMiddleware.RequestIDKey).(string)
		w.Header().Set(chiMiddleware.RequestIDHeader, reqID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
