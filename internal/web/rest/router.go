package rest

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Router struct {
	logger      zerolog.Logger
	fileService distrybute.FileService
	userService distrybute.UserService
}

func (r *Router) GetHttpHandler() http.Handler {
	router := chi.NewRouter()
	router.Post("/user/create", r.handleUserCreate)
	return router
}

type Response struct {
	StatusCode   int
	ErrorMessage string      `json:"error_message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

func (r *Router) sendResponse(w http.ResponseWriter, req *http.Request, data interface{}) {
	r.sendResponseWithCode(w, req, http.StatusOK, data)
}

func (r *Router) sendResponseWithCode(w http.ResponseWriter, req *http.Request, code int, data interface{}) {
	resp := &Response{
		StatusCode: http.StatusOK,
		Data:       data,
	}
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Err(err).Msg("could not write http response")
		return
	}
}

func (r *Router) sendNotFound(w http.ResponseWriter, req *http.Request, message string) {
	r.sendError(w, req, http.StatusNotFound, message)
}

func (r *Router) sendInternalServerError(w http.ResponseWriter, req *http.Request) {
	r.sendAutomaticError(w, req, http.StatusInternalServerError)
}

func (r *Router) sendAutomaticError(w http.ResponseWriter, req *http.Request, code int) {
	r.sendError(w, req, code, http.StatusText(code))
}

func (*Router) sendError(w http.ResponseWriter, req *http.Request, code int, message string) {
	w.WriteHeader(code)
	response := &Response{
		StatusCode:   http.StatusNotFound,
		ErrorMessage: message,
	}
	encodedResponse, err := json.Marshal(response)
	if err != nil {
		log.Err(err).Msg("could not marshal http error response")
		return
	}
	if bytesWritten, err := w.Write(encodedResponse); err != nil {
		log.Err(err).Int("bytesWritten", bytesWritten).Msg("could not write http error response")
		return
	}
}
