package web

import (
	"errors"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestRedirectUriBuild(t *testing.T) {
	t.Run("normal prefix", func(t *testing.T) {
		t.Parallel()
		redirectUri := getRequestRedirectUri("/distrybute/request/file/AFUhafuah", "/p")
		assert.Equal(t, "/distrybute/request/file/p/AFUhafuah", redirectUri)
	})
	t.Run("empty prefix", func(t *testing.T) {
		t.Parallel()
		redirectUri := getRequestRedirectUri("/distrybute/request/file/AFUhafuah", "")
		assert.Equal(t, "/distrybute/request/file/AFUhafuah", redirectUri)
	})
}

func TestBrowserRequestingCheck(t *testing.T) {
	testCases := []struct {
		name              string
		userAgentContains []string
		userAgent         string
		result            bool
	}{
		{"empty user agent contains", make([]string, 0), "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36", false},
		{"contains user agent", []string{"Mozilla/", "Opera/"}, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36", true},
		{"empty user agent contains", []string{"curl/"}, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.169 Safari/537.36", false},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			result := isBrowserRequesting(testCase.userAgentContains, testCase.userAgent)
			assert.Equal(t, testCase.result, result)
		})
	}
}

func TestDispositionHeaderRetrieval(t *testing.T) {
	testCases := []struct {
		name                  string
		contentTypesToDisplay []string
		filename              string
		contentType           string
		expected              string
	}{
		{"no whitelisted content types", make([]string, 0), "malicious.html", "text/html", `attachment; filename="malicious.html"`},
		{"content type not whitelisted", []string{"*", "text/plain", "application/octet-stream"}, "malicious.html", "text/html", `attachment; filename="malicious.html"`},
		{"content type whitelisted", []string{"*", "text/plain", "application/octet-stream"}, "paper.txt", "text/plain", `inline; filename="paper.txt"`},
	}
	for _, testCase := range testCases {
		// capture range variable
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			dispositionHeaderValue := getDispositionHeader(testCase.contentTypesToDisplay, testCase.contentType, testCase.filename)
			assert.Equal(t, testCase.expected, dispositionHeaderValue)
		})
	}
}

func TestEntryRequestErrorCheck(t *testing.T) {
	testCases := []struct {
		name       string
		err        error
		result     bool
		statusCode int
	}{
		{"error is nil", nil, false, http.StatusOK},
		{"error is of type entry not found", distrybute.ErrEntryNotFound, true, http.StatusNotFound},
		{"error is of unexpected type", errors.New("unexpected super crazy error"), true, http.StatusInternalServerError},
	}
	for _, testCase := range testCases {
		// capture range variable
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			recorder := httptest.NewRecorder()
			result := checkEntryRequestError(testCase.err, recorder)
			assert.Equal(t, testCase.result, result)
			assert.Equal(t, testCase.statusCode, recorder.Code)
		})
	}
}
