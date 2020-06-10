package web

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDispositionHeaderRetrieval(t *testing.T) {
	t.Run("no whitelisted content types", func(t *testing.T) {
		contentTypesToDisplay := make([]string, 0)
		filename := "malicious.html"
		dispositionHeaderValue := getDispositionHeader(contentTypesToDisplay, "text/html", filename)
		assert.Equal(t, `attachment; filename="malicious.html"`, dispositionHeaderValue)
	})
	t.Run("content type not whitelisted", func(t *testing.T) {
		contentTypesToDisplay := []string{"*", "text/plain", "application/octet-stream"}
		filename := "malicious.html"
		dispositionHeaderValue := getDispositionHeader(contentTypesToDisplay, "text/html", filename)
		assert.Equal(t, `attachment; filename="malicious.html"`, dispositionHeaderValue)
	})
	t.Run("content type whitelisted", func(t *testing.T) {
		contentTypesToDisplay := []string{"*", "text/plain", "application/octet-stream"}
		filename := "paper.txt"
		dispositionHeaderValue := getDispositionHeader(contentTypesToDisplay, "text/plain", filename)
		assert.Equal(t, `inline; filename="paper.txt"`, dispositionHeaderValue)
	})
}
