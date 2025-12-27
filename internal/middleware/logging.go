package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

var logger *zap.Logger // глобальный логгер (инициализируем в main.go)

func InitLogger(l *zap.Logger) {
	logger = l
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем ResponseWriter для получения статус-кода
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		logger.Info("HTTP запрос",
			zap.String("method", r.Method),
			zap.String("path", r.URL.Path),
			zap.String("remote_addr", r.RemoteAddr),
			zap.Int("status", lrw.statusCode),
			zap.Duration("duration", duration),
		)
	})
}

// Вспомогательная структура для захвата статус-кода
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
