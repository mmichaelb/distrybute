package web

import (
	"fmt"
	"github.com/go-chi/chi"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"net/http"
	"strings"
)

const (
	callReferenceParamName = "callReference"
	contentTypeHeader      = "Content-Type"
	dispositionHeader      = "Content-Disposition"
	dispositionValueFormat = "%v; filename=\"%v\""

	downloadEntryUriPrefix = "/d"
	viewEntryUriPrefix     = "/v"
)

// handleRequest represents an endpoint which can be used request a file entry
func (router *Router) handleRequest(writer http.ResponseWriter, req *http.Request) {
	// get the call reference from the url
	callReference := chi.URLParam(req, callReferenceParamName)
	// request file from storage
	entry, err := router.fileService.Request(callReference)
	// check if the file could not be found or a different error occurred
	if checkEntryRequestError(err, writer) {
		return
	}
	// make sure that the entry gets closed afterwards
	defer entry.ReadCloseSeeker.Close()
	// set content disposition header (download or embed)
	writer.Header().Set(dispositionHeader, getDispositionHeader(router.config.ContentTypesToDisplay, entry.ContentType, entry.Filename))
	// set content type header in order to simplify the serve content function
	writer.Header().Set(contentTypeHeader, entry.ContentType)
	// finally, serve the content with all the required caching parameters
	http.ServeContent(writer, req, "" /* parameter never used - see documentation */, entry.UploadDate, entry.ReadCloseSeeker)
}

// handleGeneralRequest handles a general file request and redirects the request according to the used user agent.
func (router *Router) handleGeneralRequest(writer http.ResponseWriter, req *http.Request) {
	var prefix string
	if isBrowserRequesting(router.config.BrowserUserAgentContains, req.UserAgent()) {
		prefix = viewEntryUriPrefix
	} else {
		prefix = downloadEntryUriPrefix
	}
	// get redirect url
	newUrl := getRequestRedirectUri(req.URL.RequestURI(), prefix)
	// redirect request
	http.Redirect(writer, req, newUrl, http.StatusTemporaryRedirect)
}

func getRequestRedirectUri(requestUri, prefix string) string {
	lastSlashIndex := strings.LastIndexByte(requestUri, '/')
	newUrl := requestUri[:lastSlashIndex] + prefix + requestUri[lastSlashIndex:]
	return newUrl
}

func isBrowserRequesting(userAgentContains []string, userAgent string) bool {
	if len(userAgent) == 0 {
		return false
	}
	for _, userAgentContain := range userAgentContains {
		if strings.Contains(userAgent, userAgentContain) {
			return true
		}
	}
	return false
}

func getDispositionHeader(contentTypesToDisplay []string, contentType, filename string) string {
	var dispositionType string
	for _, whitelistedContentType := range contentTypesToDisplay {
		if strings.EqualFold(whitelistedContentType, contentType) {
			dispositionType = "inline"
		}
	}
	if len(dispositionType) == 0 {
		dispositionType = "attachment"
	}
	return fmt.Sprintf(dispositionValueFormat, dispositionType, filename)
}

// checkEntryRequestError checks if the error is given if so, returns true
func checkEntryRequestError(err error, writer http.ResponseWriter) bool {
	if err != nil {
		switch err {
		case distrybute.ErrEntryNotFound:
			writer.WriteHeader(http.StatusNotFound)
			break
		default:
			// TODO log error
			writer.WriteHeader(http.StatusInternalServerError)
			break
		}
		return true
	}
	return false
}
