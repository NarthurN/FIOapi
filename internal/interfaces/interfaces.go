package interfaces

import "net/http"

type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
}

type MiddlewareLog interface {
	Log(next http.Handler) http.Handler
}
