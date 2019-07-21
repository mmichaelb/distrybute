package postgres

import (
	"github.com/google/uuid"
	"github.com/mmichaelb/gosharexserver"
)

func (manager *Manager) ListEntries(limit int, offset int, sortBy gosharexserver.FileEntrySortElem, sortOrder gosharexserver.SortSequence, uid []uuid.UUID) (err error, entries []*gosharexserver.FileEntry) {
	panic("not implemented")
}

func (manager *Manager) SearchEntries(query string, limit int, offset int, sortBy gosharexserver.FileEntrySortElem, sortOrder gosharexserver.SortSequence, uid []uuid.UUID) (err error, entries []*gosharexserver.FileEntry) {
	panic("not implemented")
}
