package postgres

import (
	"github.com/mmichaelb/gosharexserver"
	"io"
)

// Store is the Postgres based implementation of the Service interface Store function.
func (service *Service) Store(entry *gosharexserver.FileEntry, reader io.Reader) (err error) {
	panic("not implemented")
}

// Delete is the Postgres based implementation of the Service interface Delete function.
func (service *Service) Delete(entries ...gosharexserver.FileEntry) (deleted int64, err error) {
	panic("not implemented")
}
