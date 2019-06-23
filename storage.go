package gosharexserver

import (
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

const (
	// ContentTypeStatsMinimum holds the minimum of entries which must have the same content type
	// in order to be listed in the content type statistic call as a separate part.
	// TODO: change to dynamic configuration value
	ContentTypeStatsMinimum = 10
)

// FileManager holds all functions needed for a useable file manager implementation. It contains
// statistic functions and the upload and deletion of file entries. It may uses a storage helper
// implementation in order to separate the storage and metadata saving process.
type FileManager interface {
	// Store saves the entry data to the storage and returns the writer in order to write the file
	// to the storage. If something went wrong, an error (err) is returned.
	// In addition to that the entry`s ID and callReference fields are manipulated after the
	// function has finished.
	Store(entry *FileEntry, reader io.Reader) (err error)
	// Request searches for an entry by using the specified callReference. It returns an error if
	// something went wrong.
	Request(callReference string) (entry *FileEntry, err error)
	// Delete tries to search for the entries by using any of the following values set in the entry
	// instance: CallReference or DeleteReference. The entries are deleted if the search was
	// successful. It returns an error (err) if something went wrong and the total number of deleted
	// file entries.
	Delete(entries []*FileEntry) (err error, deleted int64)
	// ListEntries lists up all matched entries by searching for all entries by the given uuids and
	// returns the matched ones. It also accepts a various number of parameters to modify the search
	// results. It returns an error (err) if something went wrong.
	ListEntries(limit int, offset int, sortBy FileEntrySortElem, sortOrder SortSequence, uid []uuid.UUID) (
		err error, entries []*FileEntry)
	// SearchEntries searches for specific entries by using the parameters and returns the matched
	// ones. It returns an error (err) if something went wrong.
	SearchEntries(query string, limit int, offset int, sortBy FileEntrySortElem, sortOrder SortSequence, uid []uuid.UUID) (
		err error, entries []*FileEntry)
	// ContentTypeStatistics returns the content type statistics for the given uuids. The
	// MapBasedStatistics (stats) instance contains the content types as keys and the number of
	// matched file entries as values. For all content types falling below the
	// ContentTypeStatsMinimum the key "other" is used. It returns an error (err) if something went
	// wrong.
	ContentTypeStatistics(uid []uuid.UUID) (err error, stats *MapBasedStatistics)
	// UploadStatistics returns the upload statistics for a given period of time specified as a
	// parameter. The MapBasedStatistics (stats) instance contains the time stamps as strings
	// according to ISO 8601 in a simplified version which is YYYYMMDDThhmmssZ and UTC.
	// TODO: use the int64 unix time as a key.
	UploadStatistics(uid []uuid.UUID, period time.Duration) (err error, stats *MapBasedStatistics)
}

var (
	// ErrDuplicateStorageID indicates that the provided ID is already present in the file storage.
	ErrDuplicateStorageID = errors.New("the given ID is already present in the file storage")
	// ErrIDNotFound indicates that there is no such file with this ID in the file storage.
	ErrIDNotFound = errors.New("the given ID was not found in the file storage")
)

// Storage holds functions which are only being called by a file manager and this interface is used
// as a helper to separate the storing and managing process of files.
type Storage interface {
	// PutFile stores input in the given file storage. It accepts an io.Reader instance to allow a
	// streamable saving process. In addition to that, it needs an identical string which is used to
	// identify the object in further retrievals. If the ID is already present in the file storage,
	// ErrDuplicateStorageID is returned - if different errors occur, they are returned without wrapping.
	PutFile(id string, reader io.Reader) error
	// GetFile first searches for the given ID and if found, returns the content by handing over a
	// ReadCloseSeeker. If the given ID cannot be found in the file entry, ErrIDNotFound is
	// returned - if different errors occur, they are returned without wrapping.
	GetFile(id string) (ReadCloseSeeker, error)
	// DeleteFile tries to delete the file associated with the given ID. If the given ID cannot be
	// found in the file entry, ErrIDNotFound is returned - if different errors occur, they are
	// returned without wrapping.
	DeleteFile(id string) error
	// DeleteMultipleFiles does the same like DeleteFile but is able to delete multiple files at once.
	DeleteMultipleFiles(ids ...string) error
}

// MapBasedStatistics contains the key-based statistics needed for the statistics endpoint.
type MapBasedStatistics struct {
	// TotalEntryAmount holds the total amount of files matched.
	TotalEntryAmount int64
	// StatisticMap contains the key-based statistics - varies depending on the context.
	StatisticMap map[string]int64
}
