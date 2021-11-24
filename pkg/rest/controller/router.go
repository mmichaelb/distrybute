package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/rs/zerolog"
)

// @title distrybute API
// @version 0.0.1
// @description API documentation for the REST API of distrybute, a lightweight image upload server.

// @license.name MIT
// @license.url https://github.com/mmichaelb/distrybute/blob/master/LICENSE

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
