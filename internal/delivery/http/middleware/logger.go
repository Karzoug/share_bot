package middleware

import (
	"net/http"

	"go.uber.org/zap"
)

// Logger returns a logger handler.
func Logger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			mw := newWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(mw, r)

			logger.
				With(
					zap.String("proto", r.Proto),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
					zap.Int("status", mw.Status()),
					zap.Int("size", mw.BytesWritten())).
				Info("request")
		})
	}
}
