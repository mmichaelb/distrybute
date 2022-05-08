package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	distrybute "github.com/mmichaelb/distrybute/pkg"
	"github.com/mmichaelb/distrybute/pkg/mocks"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"strings"
	"testing"
)

var fileService *mocks.FileService
var userService *mocks.UserService
var r *router

type stringReadCloser struct {
	reader *strings.Reader
}

func (s stringReadCloser) Read(p []byte) (n int, err error) {
	return s.reader.Read(p)
}

func (s stringReadCloser) Seek(offset int64, whence int) (int64, error) {
	return s.reader.Seek(offset, whence)
}

func (s stringReadCloser) Close() error {
	return nil
}

func TestMain(m *testing.M) {
	log.Level(zerolog.DebugLevel)
	fileService = &mocks.FileService{}
	userService = &mocks.UserService{}
	r = NewRouter(log.Logger, fileService, userService)
	// hook file request endpoint
	r.Get("/v/{callReference}", r.HandleFileRequest)
	m.Run()
}

func TestRouter_HandleFileRequest(t *testing.T) {
	t.Run("entry can not be found", func(t *testing.T) {
		fileService.On("Request", "testnotfoundcr").Return(nil, distrybute.ErrEntryNotFound)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v/testnotfoundcr", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
	t.Run("unknown error leads to internal server error", func(t *testing.T) {
		fileService.On("Request", "testunknownerror").Return(nil, errors.New("some unknown error"))
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v/testunknownerror", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("entry can be requested", func(t *testing.T) {
		content := `{"object":{"a":"b","c":"d","e":"f"},"array":[1,2],"string":"Hello World"}`
		reader := &stringReadCloser{strings.NewReader(content)}
		entry := &distrybute.FileEntry{
			Filename:        "testfile.json",
			ContentType:     "application/json",
			ReadCloseSeeker: reader,
			Size:            int64(len(content)),
		}
		fileService.On("Request", "testrequest").Return(entry, nil)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/v/testrequest", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		receivedContent, err := io.ReadAll(recorder.Result().Body)
		assert.NoError(t, err)
		assert.Equal(t, content, string(receivedContent))
	})
}

func TestRouter_handleFileUpload(t *testing.T) {
	t.Run("does not accept an empty auth token", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/file", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
	t.Run("unauthorized auth tokens are not allowed", func(t *testing.T) {
		userService.On("GetUserByAuthorizationToken", "notauthorized").
			Return(false, nil, nil)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/file", nil)
		req.Header.Set("Authorization", "notauthorized")
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})
	t.Run("auth error results in internal server error", func(t *testing.T) {
		userService.On("GetUserByAuthorizationToken", "errorauthtoken").
			Return(false, nil, errors.New("some error"))
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/file", nil)
		req.Header.Set("Authorization", "errorauthtoken")
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("upload function works", func(t *testing.T) {
		testUuid, _ := uuid.Parse("c0bb684a-ecb4-4211-a31e-dc878bc001c7")
		bodyContent := `{"object":{"a":"b","c":"d","e":"f"},"array":[1,2],"string":"Hello World"}`
		testContentType := "application/test-content-type"
		testCallReference := "testcallreference"
		testDeleteReference := "testdeletereference"
		userService.On("GetUserByAuthorizationToken", "authorizedtoken").
			Return(true, &distrybute.User{ID: testUuid}, nil)
		validated := false
		fileService.On("Store", "", mock.AnythingOfType("string"),
			mock.AnythingOfType("int64"), testUuid, mock.Anything).
			Return(func(filename, contentType string, size int64, author uuid.UUID, reader io.Reader) *distrybute.FileEntry {
				assert.Empty(t, filename)
				assert.Equal(t, testContentType, contentType)
				assert.Equal(t, int64(len(bodyContent)), size)
				assert.Equal(t, testUuid, author)
				receivedContent, err := io.ReadAll(reader)
				assert.NoError(t, err)
				assert.Equal(t, []byte(bodyContent), receivedContent)
				validated = true
				return &distrybute.FileEntry{
					CallReference:   testCallReference,
					DeleteReference: testDeleteReference,
				}
			}, func(filename, contentType string, size int64, author uuid.UUID, reader io.Reader) error { return nil })
		recorder := httptest.NewRecorder()
		body, multipartWriter := prepareTestMultipart(t, bodyContent, testContentType)
		req := httptest.NewRequest(http.MethodPost, "/file", body)
		req.Header.Set("Authorization", "authorizedtoken")
		req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
		receivedContent, err := io.ReadAll(recorder.Result().Body)
		assert.True(t, validated, "could not validate upload function in file service")
		assert.NoError(t, err)
		respJsonBody := &Response{Data: &FileUploadResponse{}}
		err = json.Unmarshal(receivedContent, respJsonBody)
		assert.NoError(t, err)
		assert.Equal(t, testCallReference, respJsonBody.Data.(*FileUploadResponse).CallReference)
		assert.Equal(t, testDeleteReference, respJsonBody.Data.(*FileUploadResponse).DeleteReference)
	})
	t.Run("invalid post request is being handled normally", func(t *testing.T) {
		testUuid, _ := uuid.Parse("c0bb684a-ecb4-4211-a31e-dc878bc001c7")
		userService.On("GetUserByAuthorizationToken", "authorizedtoken").
			Return(true, &distrybute.User{ID: testUuid}, nil)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/file", nil)
		req.Header.Set("Authorization", "authorizedtoken")
		req.Header.Set("Content-Type", "invalid/contenttype")
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusBadRequest, recorder.Code)
		receivedContent, err := io.ReadAll(recorder.Result().Body)
		assert.NoError(t, err)
		respJsonBody := &Response{Data: &FileUploadResponse{}}
		err = json.Unmarshal(receivedContent, respJsonBody)
		assert.NoError(t, err)
	})
}

func prepareTestMultipart(t *testing.T, body string, contentType string) (io.Reader, *multipart.Writer) {
	var buffer bytes.Buffer
	multipartWriter := multipart.NewWriter(&buffer)
	h := make(textproto.MIMEHeader)
	h.Set("Content-Disposition", `form-data; name="file"; filename="test.json"`)
	h.Set("Content-Type", contentType)
	writer, err := multipartWriter.CreatePart(h)
	assert.NoError(t, err)
	_, err = writer.Write([]byte(body))
	assert.NoError(t, err)
	err = multipartWriter.Close()
	assert.NoError(t, err)
	return &buffer, multipartWriter
}

func TestRouter_handleFileDeletion(t *testing.T) {
	t.Run("unknown delete reference does not lead to action", func(t *testing.T) {
		fileService.On("Delete", "notfound").Return(distrybute.ErrEntryNotFound)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/file/delete/notfound", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusNotFound, recorder.Code)
	})
	t.Run("internal error leads to 500 status code", func(t *testing.T) {
		fileService.On("Delete", "servererror").Return(errors.New("some unknown error"))
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/file/delete/servererror", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	})
	t.Run("valid delete reference leads to deletion", func(t *testing.T) {
		fileService.On("Delete", "validref").Return(nil)
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/file/delete/validref", nil)
		r.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}
