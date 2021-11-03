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
// @Summary Retrieve a file by using the callReference parameter.
// @ID uploadFile
// @Accept json
// @Produce octet-stream
// @Param callReference path int true "Call Reference"
// @Success 200
// @Router /v/{callReference} [get]
func (r *router) HandleFileRequest(w http.ResponseWriter, req *http.Request) {
	writer := r.wrapResponseWriter(w)
	// retrieve file reference from request
	callReference := chi.URLParam(req, FileRequestShortIdParamName)
	// request file entry from backend
	entry, err := r.fileService.Request(callReference)
	if err == distrybute.ErrEntryNotFound {
		writer.WriteNotFoundResponse("entry not found", nil, req)
		return
	} else if err != nil {
		r.logger.Err(err).Str("callReference", callReference).Msg("could not request file entry")
		writer.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	defer func(ReadCloseSeeker distrybute.ReadCloseSeeker) {
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

// handleFileUpload handles an incoming file upload.
// @Summary Upload a file using a POST request.
// @Accept multipart/form-data
// @Produce json
// @Success 200
// @Router /api/file [post]
func (r *router) handleFileUpload(w http.ResponseWriter, req *http.Request) {

}
