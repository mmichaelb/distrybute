package postgresminio

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_PostgresMinio_UserService(t *testing.T) {
	connection, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	defer connection.Close(context.Background())
	assert.NoError(t, err)
}
