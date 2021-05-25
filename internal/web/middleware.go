package web

import (
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"net/http"
)

type LoggingMiddleware struct {
	logger zerolog.Logger
}

func (m *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ww := middleware.NewWrapResponseWriter(w, req.ProtoMajor)
	defer func() {
		status := ww.Status()
		m.logger.Info().
			Str("adr", req.RemoteAddr).
			Str("method", req.Method).
			Str("path", req.URL.RequestURI()).
			Int("status", status).
			Msg(http.StatusText(status))
	}()
}
