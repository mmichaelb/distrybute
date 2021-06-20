package postgresminio

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var connection *pgx.Conn
var minioClient *minio.Client

var testBucketName = os.Getenv("TEST_MINIO_BUCKET_NAME")

func Test_PostgresMinio_Service(t *testing.T) {
	setupPostgresConnection(t)
	setupMinioClient(t)
	defer connection.Close(context.Background())
	service := NewService(connection, minioClient, testBucketName, "")
	err := service.InitDDL()
	assert.NoError(t, err)
}

func setupPostgresConnection(t *testing.T) {
	host := os.Getenv("TEST_POSTGRES_HOST")
	port := os.Getenv("TEST_POSTGRES_PORT")
	db := os.Getenv("TEST_POSTGRES_DB")
	var err error
	connection, err = pgx.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/%s", host, port, db))
	assert.NoError(t, err, "could not establish test connection")

}

func setupMinioClient(t *testing.T) {
	host := os.Getenv("TEST_MINIO_HOST")
	port := os.Getenv("TEST_MINIO_PORT")
	var err error
	minioClient, err = minio.New(fmt.Sprintf("%s:%s", host, port), &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	assert.NoError(t, err, "could not create minio client")
	err = minioClient.MakeBucket(context.Background(), testBucketName, minio.MakeBucketOptions{})
	assert.NoError(t, err, "could not create distrybute minio test bucket")
	t.Cleanup(func() {
		err = minioClient.RemoveBucket(context.Background(), testBucketName)
		assert.NoError(t, err)
	})
}
