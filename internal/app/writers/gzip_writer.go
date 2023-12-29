package writers

import (
	"compress/gzip"
	"net/http"
)

type GzipWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

func NewGzipWriter(w http.ResponseWriter) *GzipWriter {
	return &GzipWriter{
		w:  w,
		zw: gzip.NewWriter(w),
	}
}

func (c *GzipWriter) Header() http.Header {
	return c.w.Header()
}

func (c *GzipWriter) Write(p []byte) (int, error) {
	return c.zw.Write(p)
}

func (c *GzipWriter) WriteHeader(statusCode int) {
	if statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *GzipWriter) Close() error {
	return c.zw.Close()
}
