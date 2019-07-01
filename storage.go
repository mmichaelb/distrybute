package gosharexserver

import (
	"errors"
	"io"
	"time"

	"github.com/google/uuid"
)

// Period is used when passing a period of time e.g. when retrieving statistics.
type Period string

const (
	// TODO make limits dynamic

	// PeriodHour indicates the number of uploaded files per hour (last 24h) -> return number: 24
	PeriodHour = Period("HOUR")
	// PeriodDay indicates the number of uploaded files per day (last 30 days) -> return number: 30
	PeriodDay = Period("DAY")
	// PeriodMonth indicates the number of uploaded files per month (last 12 months) -> return number: 12
	PeriodMonth = Period("MONTH")
	// PeriodYear indicates the number of uploaded files per year (last 10 years) -> return number: 10
	PeriodYear = Period("YEAR")
	// PeriodHourOfDay indicates the average number of uploaded files per hour within a day (00-23h) -> return number: 24
	PeriodHourOfDay = Period("HOUR_OF_DAY")
	// PeriodDayOfWeek indicates the average number of uploaded files per day within a week (Monday - Sunday) -> return number: 7
	PeriodDayOfWeek = Period("DAY_OF_WEEK")
)

// FileManager holds all functions needed for a useable file manager implementation. It contains
// statistic functions and the upload and deletion of file entries. It may uses a storage helper
// implementation in order to separate the storage and metadata saving process.
type FileManager interface {
	// Store saves the entry data to the storage and returns the writer in order to write the file
	// to the storage. If something went wrong, an error (err) is returned.
	// In addition to that the entry`s ID, CallReference and (maybe) DeleteReference fields are manipulated after the
	// function has finished.
	Store(entry *FileEntry, reader io.Reader) (err error)
	// Request searches for an entry by using the specified CallReference. It returns an error if
	// something went wrong.
	Request(callReference string) (entry *FileEntry, err error)
	// Delete tries to search for the entries by using any of the following values set in the entry
	// instance: ID, CallReference or DeleteReference. The entries are deleted if the search was
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
	// ResolveMIMETypeStatistic resolves the MIME type statistic for the given uuids. The
	// MIMETypeStatistic instance contains the MIME types as keys and the number of matched file
	// entries as values. The parameter uids indicates whose uploaded files should be included. It
	//  returns an error (err) if something went wrong.
	ResolveMIMETypeStatistic(uids []uuid.UUID) (err error, totalEntries int64, statistic MIMETypeStatistic)
	// ResolveUserUploadPeriodStatistic resolves the user upload statistic and sets the total number
	// of uploaded files in the UserUploadStatistic return parameter. The parameter uid indicates
	// whose uploaded files should be used. It  returns an error (err) if something went wrong.
	ResolveUserUploadPeriodStatistic(uid uuid.UUID, period Period) (err error, statistic *UserUploadPeriodStatistic)
}

// MIMETypeStatistic contains the Content-Type/MIME Type as a key and the total number of files
// using this MIME Type.
type MIMETypeStatistic map[string]int64

// UserUploadPeriodStatistic contains the values when retrieving period based upload statistics.
type UserUploadPeriodStatistic []int64

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
