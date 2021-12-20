package postgresminio

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/rs/zerolog/log"
)

func deferCloseConnFunc(conn *pgxpool.Conn) func() {
	return func() {
		if err := conn.Conn().Close(context.Background()); err != nil {
			log.Warn().Err(err).Msg("could not close postgresql pool connection")
		}
	}
}
