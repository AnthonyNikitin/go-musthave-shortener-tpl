package middlewares

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/logging"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/writers"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := logging.NewLogger()

		start := time.Now()
		responseData := writers.NewResponseData()
		lw := writers.NewLoggingResponseWriter(w, responseData)

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		logger.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	})
}
