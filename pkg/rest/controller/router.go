package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/rs/zerolog"
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
	logger      zerolog.Logger
	fileService distrybute.FileService
	userService distrybute.UserService
}

func NewRouter(logger zerolog.Logger, fileService distrybute.FileService, userService distrybute.UserService) *router {
	router := &router{
		Mux:         chi.NewRouter(),
		logger:      logger,
		fileService: fileService,
		userService: userService,
	}
	router.setupMiddlewares()
	return router
}

func (r *router) handlerFuncWithAuth(handlerFn http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		user := request.Context().Value(userContextKey).(*distrybute.User)
		if user == nil {
			r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusUnauthorized, nil, request)
			return
		}
		handlerFn.ServeHTTP(writer, request)
	}
}
