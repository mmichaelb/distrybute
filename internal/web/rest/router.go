package rest

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

type router struct {
	logger      zerolog.Logger
	fileService distrybute.FileService
	userService distrybute.UserService
}

func NewRouter(logger zerolog.Logger, fileService distrybute.FileService, userService distrybute.UserService) *router {
	return &router{logger: logger, fileService: fileService, userService: userService}
}

func (r *router) BuildHttpHandler() http.Handler {
	router := chi.NewRouter()
	middleware.RequestIDHeader = ""
	router.Use(hlog.NewHandler(r.logger))
	router.Use(middleware.CleanPath)
	router.Use(hlog.RequestIDHandler("request_id", ""))
	router.Use(r.loggingMiddleware)
	router.Use(r.recovererMiddleware)
	router.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusNotFound, nil, request)
	})
	router.MethodNotAllowed(func(writer http.ResponseWriter, request *http.Request) {
		r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusMethodNotAllowed, nil, request)
	})
	router.Post("/user/create", r.wrapStandardHttpMethod(r.handleUserCreate))
	return router
}

func (r *router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		hlog.FromRequest(request).Info().
			Str("addr", request.RemoteAddr).
			Str("userAgent", request.Header.Get("User-Agent")).
			Str("method", request.Method).
			Str("path", request.RequestURI).
			Send()
		wrappedWriter := r.wrapResponseWriter(writer)
		defer func() {
			hlog.FromRequest(request).Info().
				Int("responseCode", wrappedWriter.statusCode).Send()
		}()
		next.ServeHTTP(wrappedWriter, request)
	})
}

func (r *router) recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			recoveredValue := recover()
			if recoveredValue == nil || recoveredValue == http.ErrAbortHandler {
				return
			}
			if err, ok := recoveredValue.(error); ok {
				hlog.FromRequest(request).Err(err).Msg("recovered an error from an http handler")
			} else {
				hlog.FromRequest(request).Error().Str("recoveredValue", fmt.Sprintf("%+v", recoveredValue)).
					Msg("recovered an unknown value from an http handler")
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
