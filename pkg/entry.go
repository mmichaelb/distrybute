package distrybute

import (
	"github.com/google/uuid"
	"io"
	"time"
)

// ReadCloseSeeker combines both, io.ReaderCloser and io.Seeker into one interface. It is used within
// the FileEntry type.
// Allow access via the built in interface and implement the Read, Close and Seek methods.
type ReadCloseSeeker interface {
	io.ReadCloser
	io.Seeker
}

// FileEntry represents an uploaded file and its metadata inside the storage. It has extra fields to
// resolve the file`s content.
type FileEntry struct {
	// Id holds a unique ID which identifies the Entry inside the backend. This is mandatory because the CallReference
	// can be changed by the user.
	Id uuid.UUID
	// CallReference is a unique string identifying the entry inside the database and resolving it
	// via web.
	CallReference string
	// DeleteReference is a unique string which can be used to delete the file entry with a simple
	// GET request to allow simple deletion links.
	DeleteReference string
	// Author holds a unique UserId which can be used to identify the uploader.
	Author uuid.UUID
	// Filename is the name of the file with its extension (e.g. no-virus.exe).
	Filename string
	// ContentType is the MIME-Type of the uploaded file (e.g. image/png).
	ContentType string
	// UploadDate is the exact time of when the file was uploaded.
	UploadDate time.Time
	// Reader allows to read the file entry`s content.
	ReadCloseSeeker ReadCloseSeeker
	// Size holds the total size of the file entry`s content in bytes.
	Size int64
}
