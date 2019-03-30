package gosharexserver

// SortOrder represents general sort orders for all searches.
type SortOrder int

const (
	SortAscending SortOrder = iota
	SortDescending
	SortUnsorted
)

// FileEntrySortElem holds all values to sort by when listing up file entries.
type FileEntrySortElem int

const (
	FileEntrySortName FileEntrySortElem = iota
	FileEntrySortId
	FileEntrySortDate
	FileEntrySortFileType
	FileEntrySortSize
)
