package middleware

import (
	"context"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

type Logger struct {
	log *logrus.Logger
}

func NewLogger() *Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return &Logger{log: logger}
}

func Init(ctx context.Context, method func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fields := logrus.Fields{
			"method":     r.Method,
			"path":       r.URL.Path,
			"ip":         r.RemoteAddr,
			"user_agent": r.UserAgent(),
		}

		ctx := logrus.WithFields(fields)
		ctx.Info("Received request")
		method(w, r)

		latency := time.Since(start)
		ctx.WithFields(logrus.Fields{
			"latency": latency,
		}).Info("Request completed")
	}
}
