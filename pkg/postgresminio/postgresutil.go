package postgresminio

import (
	"github.com/jackc/pgx/v4/pgxpool"
)

func deferReleaseConnFunc(conn *pgxpool.Conn) func() {
	return func() {
		conn.Release()
	}
}
