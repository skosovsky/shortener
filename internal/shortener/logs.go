package shortener

import (
	"fmt"
	"net/http"
	"time"

	"shortener/internal/log"
)

type (
	logResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}

	responseData struct {
		statusCode int
		body       []byte
		bodySize   int
	}
)

func (l *logResponseWriter) Write(data []byte) (int, error) {
	size, err := l.ResponseWriter.Write(data)
	if err != nil {
		return 0, fmt.Errorf("error writing response: %w", err)
	}

	l.responseData.bodySize += size

	return size, nil
}

func (l *logResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.statusCode = statusCode
}

func (l *logResponseWriter) Done() {
	if l.responseData.statusCode == 0 {
		l.responseData.statusCode = http.StatusOK
	}
}

func WithLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		var respData = new(responseData)

		logWriter := logResponseWriter{
			ResponseWriter: w,
			responseData:   respData,
		}

		next.ServeHTTP(&logWriter, r)

		logWriter.Done()

		duration := time.Since(start)

		log.Info("request", //nolint:contextcheck // no ctx
			log.StringAttr("uri", r.RequestURI),
			log.StringAttr("method", r.Method),
			log.StringAttr("duration", duration.String()))

		log.Info("answer", //nolint:contextcheck // no ctx
			log.IntAttr("status", respData.statusCode),
			log.IntAttr("size", respData.bodySize))

		log.Debug("additional", //nolint:contextcheck // no ctx
			log.StringAttr("Content-Type", r.Header.Get("Content-Type")),
			log.StringAttr("Content-Encoding", r.Header.Get("Content-Encoding")))
	})
}
