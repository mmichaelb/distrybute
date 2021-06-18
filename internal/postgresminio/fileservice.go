package postgresminio

import (
	"context"
	"crypto/rand"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go"
	distrybute "github.com/mmichaelb/distrybute/internal"
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

func (s *service) initFileServiceDDL() error {
	row := s.connection.QueryRow(context.Background(), `CREATE TABLE IF NOT EXISTS distrybute.entries (
		id uuid,
		author uuid NOT NULL,
		call_reference varchar(4) NOT NULL,
		delete_reference varchar(12) NOT NULL ,
		filename varchar(256) NOT NULL,
		content_type varchar(127) NOT NULL,
		upload_date timestamptz NOT NULL,
		size bigint,
		CONSTRAINT entries_pk PRIMARY KEY (id),
		CONSTRAINT entries_fk FOREIGN KEY (author) REFERENCES distrybute.users(id),
		CONSTRAINT entries_call_reference_unique UNIQUE (call_reference),
		CONSTRAINT entries_delete_reference_unique UNIQUE (delete_reference)
	);`)
	if err := row.Scan(); !errors.Is(err, pgx.ErrNoRows) {
		log.Err(err).Msg("could not run initial file service ddl")
		return err
	}
	return nil
}

func (s *service) Store(filename, contentType string, size int64, author uuid.UUID, reader io.Reader) (entry *distrybute.FileEntry, err error) {
	tx, err := s.connection.Begin(context.Background())
	if err != nil {
		return nil, err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if !errors.Is(err, pgx.ErrTxClosed) {
			log.Err(err).Str("filename", filename).Str("contentType", contentType).Msg("could not close transaction opened in order to store a new entry")
		}
	}(tx, context.Background())
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
	_, err = s.minioClient.PutObject(s.bucketName, s.objectPrefix+id.String(), reader, size, minio.PutObjectOptions{})
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

func (s *service) Request(callReference string) (entry *distrybute.FileEntry, err error) {
	row := s.connection.QueryRow(context.Background(),
		`SELECT (id, author, delete_reference, content_type, filename, size, upload_date) FROM distrybute.entries WHERE call_reference=$1`, callReference)
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
	object, err := s.minioClient.GetObject(s.bucketName, s.objectPrefix+id.String(), minio.GetObjectOptions{})
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

func (s *service) Delete(deleteReference string) (err error) {
	row := s.connection.QueryRow(context.Background(),
		`DELETE FROM distrybute.entries WHERE delete_reference=$1 RETURNING call_reference`, deleteReference)
	var callReference string
	if err := row.Scan(&callReference); errors.Is(err, pgx.ErrNoRows) {
		return distrybute.ErrEntryNotFound
	} else if err != nil {
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
