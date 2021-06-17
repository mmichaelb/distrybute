package postgres

import (
	"github.com/jackc/pgx/v4"
)

type service struct {
	connection *pgx.Conn
}

func NewService(connection *pgx.Conn) *service {
	return &service{connection: connection}
}
