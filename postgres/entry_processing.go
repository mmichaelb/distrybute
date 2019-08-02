package postgres

import (
	"github.com/google/uuid"
	"github.com/mmichaelb/gosharexserver"
)

// ListEntries is the Postgres based implementation of the Manager interface ListEntries function.
func (manager *Manager) ListEntries(limit int, offset int, sortBy gosharexserver.FileEntrySortElem, sortOrder gosharexserver.SortSequence, uid []uuid.UUID) (entries []*gosharexserver.FileEntry, err error) {
	panic("not implemented")
}

// SearchEntries is the Postgres based implementation of the Manager interface SearchEntries function.
func (manager *Manager) SearchEntries(query string, limit int, offset int, sortBy gosharexserver.FileEntrySortElem, sortOrder gosharexserver.SortSequence, uid []uuid.UUID) (entries []*gosharexserver.FileEntry, err error) {
	panic("not implemented")
}
