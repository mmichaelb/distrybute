package gosharexserver

import (
	"github.com/google/uuid"
	"io"
	"time"
)

const (
	// MinimumEntryNumberContentTypeStats indicates the minimum of entries which must have the
	// content in order to be listed in the content type statistic call.
	// TODO: change to dynamic configuration value
	ContentTypeStatsMinimum = 10
)

// Storage holds all functions needed for a useable Storage implementation. It contains statistic
// functions and the upload and deletion of file entries.
type Storage interface {
	// Store saves the entry data to the storage and returns the writer in order to write the file
	// to the storage. If something went wrong, an error (err) is returned.
	// In addition to that the entry`s ID and callReference fields are manipulated after the
	// function has finished.
	Store(entry *FileEntry) (writer io.WriteCloser, err error)
	// Request searches for an entry by using the specified callReference. It returns an error if
	// something went wrong.
	Request(callReference string) (entry *FileEntry, err error)
	// Delete tries to search for the entries by using any of the following values set in the entry
	// instance: CallReference or DeleteReference. The entries are deleted if the search was
	// successful. It returns an error (err) if something went wrong and the total number of deleted
	// file entries).
	Delete(entries []*FileEntry) (err error, deleted int64)
	// ListEntries lists up all matched entries by searching for all entries by the given uuids and
	// returns the matched ones. It also accepts a various number of parameters to modify the search
	// results. It returns an error (err) if something went wrong.
	ListEntries(limit int, offset int, sortBy FileEntrySortElem, sortOrder SortOrder, uid []uuid.UUID) (
		err error, entries []*FileEntry)
	// SearchEntries searches for specific entries by using the parameters and returns the matched
	// ones. It returns an error (err) if something went wrong.
	SearchEntries(query string, limit int, offset int, sortBy FileEntrySortElem, sortOrder SortOrder, uid []uuid.UUID) (
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

// MapBasedStatistics contains the key-based statistics needed for the statistics endpoint.
type MapBasedStatistics struct {
	// TotalEntryAmount holds the total amount of files matched.
	TotalEntryAmount int64
	// StatisticMap contains the key-based statistics - varies depending on the context.
	StatisticMap map[string]int64
}
