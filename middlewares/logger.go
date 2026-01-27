package middlewares

import (
	"log"
	"net/http"
	"time"
)

type CustomResponseWrite struct {
	http.ResponseWriter
	statusCode int
}

func (w *CustomResponseWrite) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		cw := &CustomResponseWrite{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		next.ServeHTTP(cw, r)
		latency := time.Since(startTime)
		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.Path,
			cw.statusCode,
			latency,
		)
	})
}
