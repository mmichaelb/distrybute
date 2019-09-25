package gosharexserver

import (
	"errors"
	"io"

	"github.com/google/uuid"
)

// Period is used when passing a period of time e.g. when retrieving statistics.
type Period string

const (
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

var (
	// ErrEntryNotFound indicates that there is no such file in the file storage.
	ErrEntryNotFound = errors.New("the given entry was not found in the file storage")
)

// DeleteResponse is used within the deletion process of multiple entries. It provides information
// about all entries to delete and their current deletion state.
type DeleteResponse struct {
	// EntryID is the entry`s Id within the database.
	EntryID string
	// Err is non-nil if something failed during the deletion process - otherwise it is a nil value.
	Err error
}

// FileService holds all functions needed for a usable file service implementation. It contains
// statistic functions and the upload and deletion of file entries. It may uses a storage helper
// implementation in order to separate the storage and metadata saving process.
type FileService interface {
	// Store saves the entry data to the storage and returns the writer in order to write the file
	// to the storage. If something went wrong, an error (err) is returned.
	// In addition to that the entry`s ID, CallReference and (maybe) DeleteReference fields are manipulated after the
	// function has finished.
	Store(entry *FileEntry, reader io.Reader) (err error)
	// Request searches for an entry by using the specified CallReference. It returns an error if
	// something went wrong.
	Request(callReference string) (entry FileEntry, err error)
	// Delete deletes all entries by using the specified entry ids. It returns a responseChan so that the
	// requesting party can keep track of failed/successful deletions.
	// You should read permanently from the channel because it otherwise will block the deletion process.
	Delete(entries chan string) (responseChan chan DeleteResponse)
	// ListEntries lists up all matched entries by searching for all entries by the given uuids and
	// returns the matched ones. It also accepts a various number of parameters to modify the search
	// results. It returns an error (err) if something went wrong.
	ListEntries(limit int, offset int, sortBy FileEntrySortElem, sortOrder SortSequence, uuids ...uuid.UUID) (
		entries []FileEntry, err error)
	// SearchEntries searches for specific entries by using the parameters and returns the matched
	// ones. It returns an error (err) if something went wrong.
	SearchEntries(query string, limit int, offset int, sortBy FileEntrySortElem, sortOrder SortSequence, uuids ...uuid.UUID) (
		entries []FileEntry, err error)
	// ResolveMIMETypeStatistic resolves the MIME type statistic for the given uuids. The
	// MIMETypeStatistic instance contains the MIME types as keys and the number of matched file
	// entries as values. The parameter uuids indicates whose uploaded files should be included. It
	//  returns an error (err) if something went wrong.
	ResolveMIMETypeStatistic(uuids ...uuid.UUID) (totalEntries int64, statistic MIMETypeStatistic, err error)
	// ResolveUserUploadPeriodStatistic resolves the user upload statistic and sets the total number
	// of uploaded files in the UserUploadStatistic return parameter. The parameter uid indicates
	// whose uploaded files should be used. It  returns an error (err) if something went wrong.
	ResolveUserUploadPeriodStatistic(uuid uuid.UUID, period Period) (statistic UserUploadPeriodStatistic, err error)
}

// MIMETypeStatistic contains the Content-Type/MIME Type as a key and the total number of files
// using this MIME Type.
type MIMETypeStatistic map[string]int64

// UserUploadPeriodStatistic contains the values when retrieving period based upload statistics.
type UserUploadPeriodStatistic []int64
