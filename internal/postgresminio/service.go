package postgresminio

import (
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go/v7"
)

type service struct {
	connection   *pgx.Conn
	minioClient  *minio.Client
	bucketName   string
	objectPrefix string
}

func (s service) InitDDL() error {
	err := s.initUserDDL()
	if err != nil {
		return err
	}
	err = s.initSessionDDL()
	if err != nil {
		return err
	}
	err = s.initFileServiceDDL()
	if err != nil {
		return err
	}
	return nil
}

func NewService(connection *pgx.Conn, minioClient *minio.Client, bucketName string, objectPrefix string) *service {
	return &service{connection: connection, minioClient: minioClient, bucketName: bucketName, objectPrefix: objectPrefix}
}
