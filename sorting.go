package distrybute

// SortSequence represents general sort sequences for all searches.
type SortSequence int

const (
	// SortAscending sets the sort sequence to ascending.
	SortAscending SortSequence = iota
	// SortDescending sets the sort sequence to descending.
	SortDescending
	// SortUnsorted sets the sort sequence to unsorted.
	SortUnsorted
)

// FileEntrySortElem holds all values to sort by when listing up file entries.
type FileEntrySortElem int

const (
	// FileEntrySortName sets the value by which the files are sorted to the file name.
	FileEntrySortName FileEntrySortElem = iota
	// FileEntrySortID sets the value by which the files are sorted to the file ID.
	FileEntrySortID
	// FileEntrySortDate sets the value by which the files are sorted to the upload date.
	FileEntrySortDate
	// FileEntrySortFileType sets the value by which the files are sorted to the file type.
	FileEntrySortFileType
	// FileEntrySortSize sets the value by which the files are sorted to the file size.
	FileEntrySortSize
)
