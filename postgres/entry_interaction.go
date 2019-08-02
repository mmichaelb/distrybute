package postgres

import (
	"github.com/mmichaelb/gosharexserver"
	"io"
)

// Store is the Postgres based implementation of the Manager interface Store function.
func (manager *Manager) Store(entry *gosharexserver.FileEntry, reader io.Reader) (err error) {
	panic("not implemented")
}

// Request is the Postgres based implementation of the Manager interface Request function.
func (manager *Manager) Request(callReference string) (entry *gosharexserver.FileEntry, err error) {
	panic("not implemented")
}

// Delete is the Postgres based implementation of the Manager interface Delete function.
func (manager *Manager) Delete(entries []*gosharexserver.FileEntry) (deleted int64, err error) {
	panic("not implemented")
}
