package rest

import (
	"github.com/go-chi/chi/v5"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog/log"
	"net/http"
)

const (
	fileRequestShortIdParamName = "callReference"
)

// handleFileRequest handles an incoming file request (e.g. /v/{callReference})
func (r *router) handleFileRequest(w http.ResponseWriter, req *http.Request) {
	// retrieve file reference from request
	callReference := chi.URLParam(req, fileRequestShortIdParamName)
	// request file entry from backend
	entry, err := r.fileService.Request(callReference)
	if err == distrybute.ErrEntryNotFound {
		r.sendNotFound(w, req, "entry not found")
		return
	} else if err != nil {
		r.logger.Err(err).Str("callReference", callReference).Msg("could not request file entry")
		r.sendAutomaticError(w, req, http.StatusInternalServerError)
		return
	}
	defer func(ReadCloseSeeker distrybute.ReadCloseSeeker) {
		if err := ReadCloseSeeker.Close(); err != nil {
			log.Err(err).Str("callReference", callReference).Msg("could not close file entry")
		}
	}(entry.ReadCloseSeeker)
	// set content type from file entry
	w.Header().Set("Content-Type", entry.ContentType)
	log.Debug().Str("id", entry.Id.String()).Msg("serving file entry")
	// serve content
	http.ServeContent(w, req, entry.Filename, entry.UploadDate, entry.ReadCloseSeeker)
}
