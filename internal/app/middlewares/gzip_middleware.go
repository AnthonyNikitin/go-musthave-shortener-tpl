package middlewares

import (
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/app/writers"
	"github.com/AnthonyNikitin/go-musthave-shortener-tpl/internal/readers"
	"net/http"
	"strings"
)

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		responseWriter := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			gzipWriter := writers.NewGzipWriter(w)
			responseWriter = gzipWriter
			defer gzipWriter.Close()
		}

		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			compressReader, err := readers.NewGzipReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = compressReader
			defer compressReader.Close()
		}

		next.ServeHTTP(responseWriter, r)
	})
}
