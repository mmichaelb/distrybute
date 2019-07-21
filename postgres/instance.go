package postgres

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/mmichaelb/gosharexserver"
	"io"
	"time"

	// import postgres driver
	_ "github.com/lib/pq"
)

// Manager is the Postgres based implementation of both, the gosharexserver.FileManager and the
// gosharexserver.UserService
type Manager struct {
	// implemented interfaces
	gosharexserver.FileManager
	gosharexserver.UserService
	db      *sql.DB
	storage gosharexserver.Storage
}

// New instantiates a new instance of the Postgres based Manager.
func New(db *sql.DB, storage gosharexserver.Storage) *Manager {
	return &Manager{
		db:      db,
		storage: storage,
	}
}

// NewWithConnectionString instantiates a new instance by using the given data source name (aka
// connect url)
func NewWithConnectionString(connectionString string, storage gosharexserver.Storage) (manager *Manager, err error) {
	manager = &Manager{}
	manager.db, err = sql.Open("postgresql", connectionString)
	return
}
