package postgres

import (
	"github.com/mmichaelb/gosharexserver"
	"io"
)

// Store is the Postgres based implementation of the FileManager interface Store function.
func (manager *Manager) Store(entry *gosharexserver.FileEntry, reader io.Reader) (err error) {
	entry.
	id := 
	manager.storage.PutFile(entry., reader)
}

func (manager *Manager) Request(callReference string) (entry *gosharexserver.FileEntry, err error) {
	panic("not implemented")
}

func (manager *Manager) Delete(entries []*gosharexserver.FileEntry) (err error, deleted int64) {
	panic("not implemented")
}
