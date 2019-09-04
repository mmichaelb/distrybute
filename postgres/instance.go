package postgres

import (
	"database/sql"
	"github.com/mmichaelb/gosharexserver"

	// import postgres driver
	_ "github.com/lib/pq"
)

// Service is the Postgres based implementation of three, the gosharexserver.FileService, the
// gosharexserver.UserService and gosharexserver.SessionService.
type Service struct {
	// implemented interfaces
	gosharexserver.FileService
	gosharexserver.UserService
	gosharexserver.SessionService
	db      *sql.DB
	storage gosharexserver.Storage
}

// New instantiates a new instance of the Postgres based Service.
func New(db *sql.DB, storage gosharexserver.Storage) *Service {
	return &Service{
		db:      db,
		storage: storage,
	}
}

// NewWithConnectionString instantiates a new instance by using the given data source name (aka
// connect url).
func NewWithConnectionString(connectionString string, storage gosharexserver.Storage) (service *Service, err error) {
	service = &Service{
		storage:storage,
	}
	service.db, err = sql.Open("postgresql", connectionString)
	return
}
