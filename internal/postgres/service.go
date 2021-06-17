package postgres

import (
	"github.com/jackc/pgx"
)

type service struct {
	connection *pgx.Conn
}

func NewService(connection *pgx.Conn) *service {
	return &service{connection: connection}
}
