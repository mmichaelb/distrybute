package distrybute

import (
	"errors"
	"github.com/google/uuid"
	"io"
)

var (
	// ErrEntryNotFound indicates that there is no such file in the file storage.
	ErrEntryNotFound = errors.New("the given entry was not found in the file storage")
)

// FileService holds all functions needed for a usable file service implementation.
type FileService interface {
	// Store saves the entry data to the storage. If something went wrong, an error is returned.
	Store(filename, contentType string, size int64, author uuid.UUID, reader io.Reader) (entry *FileEntry, err error)
	// Request searches for an entry by using the specified CallReference. It returns an error if something goes wrong.
	Request(callReference string) (entry *FileEntry, err error)
	// Delete deletes an entry using the provided delete reference. It returns an error if something goes wrong.
	Delete(deleteReference string) (err error)
}
