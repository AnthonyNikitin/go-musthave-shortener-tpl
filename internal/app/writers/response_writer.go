package writers

import "net/http"

type ResponseWriter interface {
	Header() http.Header
	Write([]byte) (int, error)
	WriteHeader(statusCode int)
}

type ResponseData struct {
	Status int
	Size   int
}

func NewResponseData() *ResponseData {
	return &ResponseData{
		Status: 0,
		Size:   0,
	}
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	ResponseData *ResponseData
}

func NewLoggingResponseWriter(w http.ResponseWriter, r *ResponseData) LoggingResponseWriter {
	return LoggingResponseWriter{
		ResponseWriter: w,
		ResponseData:   r,
	}
}

func (r *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.ResponseData.Size += size
	return size, err
}

func (r *LoggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.ResponseData.Status = statusCode
}
