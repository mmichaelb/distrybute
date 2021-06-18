package postgresminio

import (
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go"
)

type service struct {
	connection   *pgx.Conn
	minioClient  *minio.Client
	bucketName   string
	objectPrefix string
}

func NewService(connection *pgx.Conn) *service {
	return &service{connection: connection}
}
