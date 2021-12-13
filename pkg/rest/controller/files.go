package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/rs/zerolog/hlog"
	"mime/multipart"
	"net/http"
)

const (
	AuthorizationHeaderKey      = "Authorization"
	FileRequestShortIdParamName = "callReference"
	maximumMemoryBytes          = 1 << 20 // 1 MB maximum in memory
	multipartFormName           = "file"
)

// HandleFileRequest handles an incoming file request (e.g. /v/{callReference})
// @Router    /v/{callReference} [get]
// @ID        retrieveFile
// @Tags      files
// @Summary   Retrieve a file by using the callReference parameter.
// @Param     callReference  path  int  true  "Call Reference"
// @Produce   octet-stream,json
// @Success   200
// @Response  default  {object}  controller.Response
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
	hlog.FromRequest(req).Info().Str("id", entry.Id.String()).Msg("serving file entry")
	// serve content
	http.ServeContent(writer, req, entry.Filename, entry.UploadDate, entry.ReadCloseSeeker)
}

// handleFileUpload handles an incoming file upload.
// @Router    /api/file [post]
// @Security  ApiKeyAuth
// @ID        uploadFile
// @Tags      files
// @Summary   Upload a file using a POST request.
// @Accept    multipart/form-data
// @Param     file  formData  string  true  "Contains the file content which should be uploaded"  binary
// @Produce   json
// @success   200      {object}  controller.Response{data=controller.FileUploadResponse}  "The response which contains the callReference"
// @Response  default  {object}  controller.Response
func (r *router) handleFileUpload(w *responseWriter, req *http.Request) {
	token := req.Header.Get(AuthorizationHeaderKey)
	if token == "" {
		w.WriteAutomaticErrorResponse(http.StatusUnauthorized, nil, req)
		return
	}
	ok, user, err := r.userService.GetUserByAuthorizationToken(token)
	if err != nil {
		hlog.FromRequest(req).Err(err).Str("tokenHeader", token).Msg("could not get user by auth token")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	if !ok {
		w.WriteAutomaticErrorResponse(http.StatusUnauthorized, nil, req)
		return
	}
	// parse multipart form file and if something goes wrong return an internal server error response code
	if err = req.ParseMultipartForm(maximumMemoryBytes); err != nil {
		hlog.FromRequest(req).Warn().Err(err).Msg("could not parse multipart form")
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, nil, req)
		return
	}
	var file multipart.File
	// parse filename and mime type from multipart header
	var multipartFileHeader *multipart.FileHeader
	if file, multipartFileHeader, err = req.FormFile(multipartFormName); err != nil {
		hlog.FromRequest(req).Warn().Err(err).Msg("could not resolve multipart form details")
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, nil, req)
		return
	}
	mimeType := multipartFileHeader.Header.Get("Content-Type")
	entry, err := r.fileService.Store("", mimeType, multipartFileHeader.Size, user.ID, file)
	if err != nil {
		hlog.FromRequest(req).Err(err).Msg("could not store file entry")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	hlog.FromRequest(req).Info().
		Str("id", entry.Id.String()).
		Str("callReference", entry.CallReference).
		Int64("size", multipartFileHeader.Size).
		Msg("created new entry")
	// send json response
	w.WriteSuccessfulResponse(&FileUploadResponse{CallReference: entry.CallReference, DeleteReference: entry.DeleteReference}, req)
}

// FileUploadResponse is used to return information about an uploaded file.
type FileUploadResponse struct {
	CallReference   string `json:"callReference"`
	DeleteReference string `json:"deleteReference"`
}

// handleFileDeletion handles an incoming file deletion request.
// @Router    /api/file/delete/{deleteReference} [get]
// @ID        deleteFile
// @Tags      files
// @Summary   Deletes a specific file using the provided delete reference.
// @Param     deleteReference  path  int  true  "Call Reference"
// @Produce   json
// @Success   200      {object}  controller.Response
// @Response  default  {object}  controller.Response
func (r *router) handleFileDeletion(w *responseWriter, req *http.Request) {
	deleteReference := chi.URLParam(req, "deleteReference")
	if deleteReference == "" {
		w.WriteAutomaticErrorResponse(http.StatusBadRequest, nil, req)
		return
	}
	err := r.fileService.Delete(deleteReference)
	if err == distrybute.ErrEntryNotFound {
		w.WriteNotFoundResponse("no entry associated with the given delete reference", nil, req)
		return
	} else if err != nil {
		hlog.FromRequest(req).Err(err).Msg("could not delete entry using delete reference")
		w.WriteAutomaticErrorResponse(http.StatusInternalServerError, nil, req)
		return
	}
	hlog.FromRequest(req).Info().Str("deleteReference", deleteReference).Msg("file entry deleted")
	w.WriteSuccessfulResponse(nil, req)
}
