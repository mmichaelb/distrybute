package rest

import (
	"github.com/go-chi/chi/v5"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog"
	"net/http"
)

const (
	fileRequestShortIdParamName = "callReference"
)

// handleFileRequest handles an incoming file request (e.g. /v/{callReference})
func (r *router) handleFileRequest(w responseWriter, req *http.Request) {
	// retrieve file reference from request
	callReference := chi.URLParam(req, fileRequestShortIdParamName)
	// request file entry from backend
	entry, err := r.fileService.Request(callReference)
	if err == distrybute.ErrEntryNotFound {
		w.WriteNotFoundResponse("entry not found", nil, req)
		return
	} else if err != nil {
		r.logger.Err(err).Str("callReference", callReference).Msg("could not request file entry")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	defer func(ReadCloseSeeker distrybute.ReadCloseSeeker) {
		if err := ReadCloseSeeker.Close(); err != nil {
			r.log(zerolog.ErrorLevel, req).Err(err).Str("callReference", callReference).Msg("could not close file entry")
		}
	}(entry.ReadCloseSeeker)
	// set content type from file entry
	w.Header().Set("Content-Type", entry.ContentType)
	r.log(zerolog.DebugLevel, req).Str("id", entry.Id.String()).Msg("serving file entry")
	// serve content
	http.ServeContent(&w, req, entry.Filename, entry.UploadDate, entry.ReadCloseSeeker)
}
