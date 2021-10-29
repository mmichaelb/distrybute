package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/rs/zerolog/hlog"
	"net/http"
)

const (
	FileRequestShortIdParamName = "callReference"
)

// HandleFileRequest handles an incoming file request (e.g. /v/{callReference})
func (r *router) HandleFileRequest(w http.ResponseWriter, req *http.Request) {
	writer := r.wrapResponseWriter(w)
	// retrieve file reference from request
	callReference := chi.URLParam(req, FileRequestShortIdParamName)
	// request file entry from backend
	entry, err := r.fileService.Request(callReference)
	if err == pkg.ErrEntryNotFound {
		writer.WriteNotFoundResponse("entry not found", nil, req)
		return
	} else if err != nil {
		r.logger.Err(err).Str("callReference", callReference).Msg("could not request file entry")
		writer.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	defer func(ReadCloseSeeker pkg.ReadCloseSeeker) {
		if err := ReadCloseSeeker.Close(); err != nil {
			hlog.FromRequest(req).Err(err).Str("callReference", callReference).Msg("could not close file entry")
		}
	}(entry.ReadCloseSeeker)
	// set content type from file entry
	w.Header().Set("Content-Type", entry.ContentType)
	hlog.FromRequest(req).Debug().Str("id", entry.Id.String()).Msg("serving file entry")
	// serve content
	http.ServeContent(writer, req, entry.Filename, entry.UploadDate, entry.ReadCloseSeeker)
}