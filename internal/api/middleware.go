package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	l "github.com/deesel/wol/internal/logger"
)

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newStatusResponseWriter(w http.ResponseWriter) *statusResponseWriter {
	return &statusResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
	}
}

func (sw *statusResponseWriter) WriteHeader(statusCode int) {
	sw.statusCode = statusCode
	sw.ResponseWriter.WriteHeader(statusCode)
}

func accessLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := newStatusResponseWriter(w)

		defer func() {
			l.Get().Infow("HTTP request",
				"host", r.Host,
				"duration", fmt.Sprint(time.Since(start)),
				"status", sw.statusCode,
				"method", r.Method,
				"path", r.URL.Path,
				"query", r.URL.RawQuery,
			)
		}()

		next.ServeHTTP(sw, r)
	})
}

func jsonContent(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func handleError(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				var message string
				var statusCode int

				switch e := e.(type) {
				case error:
					statusCode = getStatusCode(e)
					message = e.Error()
				case string:
					message = e
					statusCode = http.StatusInternalServerError
				default:
					message = ErrUnknown.Error()
					statusCode = http.StatusInternalServerError
				}

				w.WriteHeader(statusCode)
				l.Get().Error(message)
				json.NewEncoder(w).Encode(map[string]interface{}{"error": message})
			}
		}()
		next.ServeHTTP(w, r)
	})
}
