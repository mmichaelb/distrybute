package postgresminio

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func Test_PostgresMinio_UserService(t *testing.T) {
	host := os.Getenv("TEST_POSTGRES_HOST")
	port := os.Getenv("TEST_POSTGRES_PORT")
	connection, err := pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/postgres", host, port))
	defer connection.Close(context.Background())
	assert.NoError(t, err)
}
