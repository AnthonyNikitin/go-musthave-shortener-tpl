package middlewares

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/writers"
	"go.uber.org/zap"
	"net/http"
	"time"
)

var sugar zap.SugaredLogger

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger, err := zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()

		sugar = *logger.Sugar()

		start := time.Now()
		responseData := writers.NewResponseData()
		lw := writers.NewLoggingResponseWriter(w, responseData)

		next.ServeHTTP(&lw, r)

		duration := time.Since(start)

		sugar.Infoln(
			"uri", r.RequestURI,
			"method", r.Method,
			"status", responseData.Status,
			"duration", duration,
			"size", responseData.Size,
		)
	})
}
