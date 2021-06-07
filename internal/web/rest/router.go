package rest

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog"
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
	router.Use(middleware.RequestID)
	router.Use(middleware.CleanPath)
	router.Use(r.recovererMiddleware)
	router.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusNotFound, request)
	})
	router.MethodNotAllowed(func(writer http.ResponseWriter, request *http.Request) {
		wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusMethodNotAllowed, request)
	})
	router.Post("/user/create", wrapStandardHttpMethod(r.handleUserCreate))
	return router
}

func (r *router) log(level zerolog.Level, request *http.Request) *zerolog.Event {
	return r.logger.WithLevel(level).
		Str("addr", request.RemoteAddr).
		Str("userAgent", request.Header.Get("User-Agent")).
		Str("requestId", middleware.GetReqID(request.Context()))
}

func (r *router) recovererMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			recoveredValue := recover()
			if recoveredValue == nil || recoveredValue == http.ErrAbortHandler {
				return
			}
			logEntry := r.log(zerolog.ErrorLevel, request)
			if err, ok := recoveredValue.(error); ok {
				logEntry.Err(err).Msg("recovered an error from an http handler")
			} else {
				logEntry.Str("recoveredValue", fmt.Sprintf("%+v", recoveredValue)).
					Msg("recovered an unknown value from an http handler")
			}
		}()
		next.ServeHTTP(writer, request)
	})
}
