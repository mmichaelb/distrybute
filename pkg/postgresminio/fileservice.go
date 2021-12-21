package postgresminio

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go/v7"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/rs/zerolog/log"
	"io"
	"math/big"
	"strings"
	"time"
)

var entryIdChars = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

const (
	callReferenceLength   = 4
	deleteReferenceLength = 12
)

func (s *Service) Store(filename, contentType string, size int64, author uuid.UUID, reader io.Reader) (entry *distrybute.FileEntry, err error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer deferReleaseConnFunc(conn)()
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func() {
		err := tx.Rollback(context.Background())
		if !errors.Is(err, pgx.ErrTxClosed) && err != nil {
			log.Err(err).Str("filename", filename).Str("contentType", contentType).Msg("could not rollback transaction opened in order to store a new entry")
		}
	}()
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	callReference, err := generateEntryId(callReferenceLength)
	if err != nil {
		return nil, err
	}
	deleteReference, err := generateEntryId(deleteReferenceLength)
	if err != nil {
		return nil, err
	}
	uploadDate := time.Now()
	row := tx.QueryRow(context.Background(),
		`INSERT INTO distrybute.entries (id, author, call_reference, delete_reference, filename, content_type, upload_date, size)
 VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`, id, author, callReference, deleteReference, filename, contentType, uploadDate, size)
	if err := row.Scan(); !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	_, err = s.minioClient.PutObject(context.Background(), s.bucketName, s.objectPrefix+id.String(), reader, size, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return nil, err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return nil, err
	}
	entry = &distrybute.FileEntry{
		Id:              id,
		CallReference:   callReference,
		DeleteReference: deleteReference,
		Author:          author,
		Filename:        filename,
		ContentType:     contentType,
		UploadDate:      uploadDate,
		Size:            size,
	}
	return entry, nil
}

func (s *Service) Request(callReference string) (entry *distrybute.FileEntry, err error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer deferReleaseConnFunc(conn)()
	row := conn.QueryRow(context.Background(),
		`SELECT id, author, delete_reference, content_type, filename, size, upload_date FROM distrybute.entries WHERE call_reference=$1`, callReference)
	var id, author uuid.UUID
	var deleteReference, contentType, filename string
	var size int64
	var uploadDate time.Time
	if err := row.Scan(&id, &author, &deleteReference, &contentType, &filename, &size, &uploadDate); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, distrybute.ErrEntryNotFound
		} else if err != nil {
			return nil, err
		}
	}
	object, err := s.minioClient.GetObject(context.Background(), s.bucketName, s.objectPrefix+id.String(), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	entry = &distrybute.FileEntry{
		Id:              id,
		CallReference:   callReference,
		DeleteReference: deleteReference,
		Author:          author,
		Filename:        filename,
		ContentType:     contentType,
		UploadDate:      uploadDate,
		ReadCloseSeeker: object,
		Size:            size,
	}
	return entry, nil
}

func (s *Service) Delete(deleteReference string) (err error) {
	conn, err := s.pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer deferReleaseConnFunc(conn)()
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer func() {
		err := tx.Rollback(context.Background())
		if !errors.Is(err, pgx.ErrTxClosed) && err != nil {
			log.Err(err).Str("deleteReference", deleteReference).Msg("could not rollback transaction opened in order to delete an entry")
		}
	}()
	row := tx.QueryRow(context.Background(),
		`DELETE FROM distrybute.entries WHERE delete_reference=$1 RETURNING id`, deleteReference)
	var id uuid.UUID
	if err := row.Scan(&id); errors.Is(err, pgx.ErrNoRows) {
		return distrybute.ErrEntryNotFound
	} else if err != nil {
		return err
	}
	err = s.minioClient.RemoveObject(context.Background(), s.bucketName, s.objectPrefix+id.String(), minio.RemoveObjectOptions{})
	if err != nil {
		return err
	}
	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

func generateEntryId(length int) (string, error) {
	var id strings.Builder
	for i := 0; i < length; i++ {
		randIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(entryIdChars))))
		if err != nil {
			return "", err
		}
		if _, err = id.WriteRune(entryIdChars[randIndex.Int64()]); err != nil {
			return "", err
		}
	}
	return id.String(), nil
}
