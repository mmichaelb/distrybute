package rest

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
	"net/http"
)

type Response struct {
	StatusCode   int         `json:"status_code"`
	ErrorMessage string      `json:"error_message,omitempty"`
	Data         interface{} `json:"data,omitempty"`
}

type responseWriter struct {
	writer     http.ResponseWriter
	statusCode int
}

func wrapResponseWriter(writer http.ResponseWriter) *responseWriter {
	return &responseWriter{
		writer: writer,
	}
}

type HandlerFunc func(*responseWriter, *http.Request)

func wrapStandardHttpMethod(handlerFunc HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		handlerFunc(wrapResponseWriter(writer), request)
	}
}

func (writer responseWriter) Header() http.Header {
	return writer.writer.Header()
}

func (writer *responseWriter) Write(bytes []byte) (int, error) {
	if writer.statusCode == 0 {
		writer.statusCode = http.StatusOK
	}
	return writer.writer.Write(bytes)
}

func (writer responseWriter) WriteHeader(statusCode int) {
	writer.writer.WriteHeader(statusCode)
	writer.statusCode = statusCode
}

func (writer responseWriter) WriteResponse(statusCode int, errorMessage string, data interface{}, r *http.Request) {
	writer.WriteHeader(statusCode)
	resp := &Response{
		StatusCode:   statusCode,
		ErrorMessage: errorMessage,
	}
	if data != nil {
		resp.Data = data
	}
	if err := json.NewEncoder(writer.writer).Encode(resp); err != nil {
		log.Err(err).Msg("could not write http response") // TODO include requesting address etc
	}
}

func (writer responseWriter) WriteSuccessfulResponse(data interface{}, r *http.Request) {
	writer.WriteResponse(http.StatusOK, "", data, r)
}

func (writer responseWriter) WriteNotFoundResponse(message string, data interface{}, r *http.Request) {
	writer.WriteResponse(http.StatusOK, message, data, r)
}

func (writer responseWriter) WriteAutomaticErrorResponse(statusCode int, r *http.Request) {
	writer.WriteResponse(statusCode, http.StatusText(statusCode), nil, r)
}
