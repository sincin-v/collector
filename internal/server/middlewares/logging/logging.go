package logging

import (
	"net/http"
	"time"

	"github.com/sincin-v/collector/internal/logger"
)

type (
	responseData struct {
		status int
		size   int
	}

	newResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func (r *newResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func (r *newResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}

func LoggerMiddleware(h http.Handler) http.Handler {
	logFunction := func(w http.ResponseWriter, r *http.Request) {

		url := r.URL.Path
		method := r.Method

		responseData := &responseData{
			status: 0,
			size:   0,
		}

		newWriter := newResponseWriter{
			responseData:   responseData,
			ResponseWriter: w,
		}

		startTime := time.Now()
		h.ServeHTTP(&newWriter, r)
		duration := time.Since(startTime)

		logger.Log.Infoln(
			"[Server] Input Request:",
			"URL:", url,
			"Method", method,
			"Duration", duration,
			"Status", responseData.status,
			"Size", responseData.size,
		)
	}

	return http.HandlerFunc(logFunction)
}
