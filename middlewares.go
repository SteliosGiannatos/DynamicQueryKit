package dynamicquerykit

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Logging Provides logging middleware for the app
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		slog.LogAttrs(context.Background(), slog.LevelInfo, "request", slog.String("method", r.Method), slog.String("path", fmt.Sprintf("%s", r.URL.Path)), slog.Any("response", time.Since(start)))
	})
}

// Cors adds CORS to all routes in the app
func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		next.ServeHTTP(w, r)
	})
}

// Middleware takes in a middleware
type Middleware func(http.Handler) http.Handler

// CreateStack provides a simple way to stack middlewares
func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(xs) - 1; i >= 0; i-- {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}
