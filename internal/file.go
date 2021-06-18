package distrybute

import (
	"errors"
	"io"
)

var (
	// ErrEntryNotFound indicates that there is no such file in the file storage.
	ErrEntryNotFound = errors.New("the given entry was not found in the file storage")
)

// FileService holds all functions needed for a usable file service implementation. It contains
// statistic functions and the upload and deletion of file entries. It may uses a storage helper
// implementation in order to separate the storage and metadata saving process.
type FileService interface {
	// Store saves the entry data to the storage and returns the writer in order to write the file
	// to the storage. If something went wrong, an error (err) is returned.
	// In addition to that the entry`s ID, CallReference and (maybe) DeleteReference fields are manipulated after the
	// function has finished.
	Store(filename, contentType string, size int64, reader io.Reader) (entry *FileEntry, err error)
	// Request searches for an entry by using the specified CallReference. It returns an error if
	// something went wrong.
	Request(callReference string) (entry *FileEntry, err error)
	// Delete deletes all entries by using the specified entry ids. It returns a responseChan so that the
	// requesting party can keep track of failed/successful deletions.
	// You should read permanently from the channel because it otherwise will block the deletion process.
	Delete(deleteReference string) (err error)
}
