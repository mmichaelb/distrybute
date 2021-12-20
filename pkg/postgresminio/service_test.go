package postgresminio

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var pool *pgxpool.Pool
var minioClient *minio.Client

var testBucketName = os.Getenv("TEST_MINIO_BUCKET_NAME")

func Test_PostgresMinio_Service(t *testing.T) {
	if os.Getenv("POSTGRES_MINIO_INTEGRATION_TEST") == "" {
		t.Skip("skipping postgres minio integration test because env `POSTGRES_MINIO_INTEGRATION_TEST` is not set")
		return
	}
	setupPostgresConnection(t)
	setupMinioClient(t)
	service := NewService(pool, minioClient, testBucketName, "")
	err := service.Init()
	assert.NoError(t, err)
	t.Run("user Service", userServiceIntegrationTest(service))
	t.Run("file Service", fileServiceIntegrationTest(service, service))
}

func setupPostgresConnection(t *testing.T) {
	host := os.Getenv("TEST_POSTGRES_HOST")
	port := os.Getenv("TEST_POSTGRES_PORT")
	db := os.Getenv("TEST_POSTGRES_DB")
	var err error
	pool, err = pgxpool.Connect(context.Background(), fmt.Sprintf("postgres://postgres:postgres@%s:%s/%s", host, port, db))
	assert.NoError(t, err, "could not establish test connection")
	conn, err := pool.Acquire(context.Background())
	assert.NoError(t, err, "could not acquire a new database connection from postgresql pool")
	err = conn.QueryRow(context.Background(), "CREATE SCHEMA IF NOT EXISTS distrybute").Scan()
	assert.ErrorIs(t, err, pgx.ErrNoRows, "could not create distrybute schema")
	t.Cleanup(func() {
		err = conn.QueryRow(context.Background(), "DROP SCHEMA distrybute CASCADE").Scan()
		assert.ErrorIs(t, err, pgx.ErrNoRows, "could not delete distrybute schema")
		err = conn.Conn().Close(context.Background())
		assert.NoError(t, err, "could not close postgres connection")
	})
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
	t.Cleanup(func() {
		objectInfoChan := minioClient.ListObjects(context.Background(), testBucketName, minio.ListObjectsOptions{})
		removeObjErrChan := minioClient.RemoveObjects(context.Background(), testBucketName, objectInfoChan, minio.RemoveObjectsOptions{})
		for {
			removeObjErr, ok := <-removeObjErrChan
			if !ok {
				break
			}
			assert.NoError(t, removeObjErr.Err, "could not remove bucket object")
		}
		err = minioClient.RemoveBucket(context.Background(), testBucketName)
		assert.NoError(t, err)
	})
}
