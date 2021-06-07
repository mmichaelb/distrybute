package rest

import (
	"github.com/go-chi/chi/v5"
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
	router.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusNotFound, request)
	})
	router.MethodNotAllowed(func(writer http.ResponseWriter, request *http.Request) {
		wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusMethodNotAllowed, request)
	})
	router.Post("/user/create", wrapStandardHttpMethod(r.handleUserCreate))
	return router
}
