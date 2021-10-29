package rest

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

// @title distrybute API
// @version 0.0.1
// @description API documentation for the REST API of distrybute, a lightweight image upload server.

// @license.name MIT
// @license.url https://github.com/mmichaelb/distrybute/blob/master/LICENSE

var userContextKey = &struct{}{}

type router struct {
	*chi.Mux
	logger         zerolog.Logger
	fileService    pkg.FileService
	userService    pkg.UserService
	sessionService pkg.SessionService
}

func NewRouter(logger zerolog.Logger, fileService pkg.FileService, userService pkg.UserService, sessionService pkg.SessionService) *router {
	router := &router{
		Mux:            chi.NewRouter(),
		logger:         logger,
		fileService:    fileService,
		userService:    userService,
		sessionService: sessionService,
	}
	router.setupMiddlewares()
	router.Post("/login", router.handleUserLogin)
	router.Post("/logout", router.handlerFuncWithAuth(router.HandleUserLogout))
	router.Post("/user/create", router.handlerFuncWithAuth(router.handleUserCreate))
	router.Get("/user/getAuthToken", router.handlerFuncWithAuth(router.handleUserRetrieveAuthToken))
	return router
}

func (r *router) setupMiddlewares() {
	middleware.RequestIDHeader = ""
	r.Use(hlog.NewHandler(r.logger))
	r.Use(middleware.CleanPath)
	r.Use(hlog.RequestIDHandler("request_id", ""))
	r.Use(r.loggingMiddleware)
	r.Use(r.recovererMiddleware)
	r.Use(r.authenticationMiddleware)
	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusNotFound, nil, request)
	})
	r.MethodNotAllowed(func(writer http.ResponseWriter, request *http.Request) {
		r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusMethodNotAllowed, nil, request)
	})
}

func (r *router) handlerFuncWithAuth(handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user := request.Context().Value(userContextKey).(*pkg.User)
		if user == nil {
			r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusUnauthorized, nil, request)
			return
		}
		handlerFn.ServeHTTP(writer, request)
	}
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

func (r *router) authenticationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var err error
		ok, user, err := r.sessionService.ValidateUserSession(request)
		if ok {
			request = request.WithContext(
				context.WithValue(request.Context(), userContextKey, user))
		}
		if err != nil {
			hlog.FromRequest(request).Err(err).Msg("could not validate user session")
			r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, request)
			return
		}
		next.ServeHTTP(writer, request)
	})
}
