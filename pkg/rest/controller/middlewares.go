package controller

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/hlog"
	"net"
	"net/http"
)

func (r *router) setupMiddlewares() {
	middleware.RequestIDHeader = ""
	r.Use(hlog.NewHandler(r.logger))
	r.Use(middleware.CleanPath)
	r.Use(hlog.RequestIDHandler("request_id", ""))
	r.Use(r.loggingMiddleware)
	r.Use(r.recovererMiddleware)
	r.NotFound(func(writer http.ResponseWriter, request *http.Request) {
		r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusNotFound, nil, request)
	})
	r.MethodNotAllowed(func(writer http.ResponseWriter, request *http.Request) {
		r.wrapResponseWriter(writer).WriteAutomaticErrorResponse(http.StatusMethodNotAllowed, nil, request)
	})
}

func (r *router) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		hlog.FromRequest(request).Info().
			Str("addr", request.RemoteAddr).
			Interface("headers", request.Header).
			Str("method", request.Method).
			Str("path", request.RequestURI).
			Msg("request.incoming")
		wrappedWriter := r.wrapResponseWriter(writer)
		go func() {
			rDns, err := net.LookupAddr(request.RemoteAddr)
			if err != nil {
				hlog.FromRequest(request).Err(err).Msg("could not look for reverse dns entry")
				return
			}
			hlog.FromRequest(request).Info().Interface("reverseDns", rDns).Msg("resolved reverse dns entries")
		}()
		defer func() {
			hlog.FromRequest(request).Info().
				Int("responseCode", wrappedWriter.statusCode).Msg("request.result")
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
