package postgresminio

import (
	"github.com/google/uuid"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"testing"
)

func fileServiceIntegrationTest(fileService distrybute.FileService, userService distrybute.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := userService.CreateNewUser("fileservice-test-user", []byte("Sommer2019"))
		assert.NoError(t, err, "could not create file service test")
		contentType := "text/plain"
		contentString := "some file content"
		content := strings.NewReader(contentString)
		size := int64(len(contentString))
		t.Run("file can be stored, retrieved and deleted", func(t *testing.T) {
			filename := "testfile.txt"
			entry, err := fileService.Store(filename, contentType, size, user.ID, content)
			assert.NoError(t, err, "entry could not be stored")
			assert.Equal(t, user.ID, entry.Author)
			assert.NotEqual(t, uuid.UUID{}, entry.Id)
			assert.Nil(t, entry.ReadCloseSeeker)
			assert.Equal(t, filename, entry.Filename)
			assert.Equal(t, size, entry.Size)
			assert.NotEmpty(t, entry.CallReference)
			assert.NotEmpty(t, entry.DeleteReference)
			assert.Equal(t, contentType, entry.ContentType)
			assert.NotEmpty(t, entry.UploadDate.Unix())
			retrievedEntry, err := fileService.Request(entry.CallReference)
			assert.NoError(t, err, "entry could not be retrieved")
			assertEntryComparison(t, entry, retrievedEntry)
			assert.NotNil(t, retrievedEntry.ReadCloseSeeker)
			contentRead, err := ioutil.ReadAll(retrievedEntry.ReadCloseSeeker)
			assert.NoError(t, err, "could not read content of entry")
			assert.Equal(t, []byte(contentString), contentRead, "returned entry content is not equal")
			err = fileService.Delete(entry.DeleteReference)
			assert.NoError(t, err, "an error occurred while deleting the test entry")
		})
		t.Run("files with duplicate names can be stored", func(t *testing.T) {
			filename := "duplicatefile.txt"
			_, err := fileService.Store(filename, contentType, size, user.ID, content)
			assert.NoError(t, err, "could not store first duplicate entry")
			_, err = fileService.Store(filename, contentType, size, user.ID, content)
			assert.NoError(t, err, "could not store second duplicate entry")
		})
	}
}

func assertEntryComparison(t *testing.T, expected *distrybute.FileEntry, actual *distrybute.FileEntry) {
	assert.Equal(t, expected.Id, actual.Id)
	assert.Equal(t, expected.Author, actual.Author)
	assert.Equal(t, expected.UploadDate.Unix(), actual.UploadDate.Unix())
	assert.Equal(t, expected.Size, actual.Size)
	assert.Equal(t, expected.ContentType, actual.ContentType)
	assert.Equal(t, expected.CallReference, actual.CallReference)
	assert.Equal(t, expected.DeleteReference, actual.DeleteReference)
	assert.Equal(t, expected.Filename, actual.Filename)
}
