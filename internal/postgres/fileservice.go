package postgres

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/mmichaelb/distrybute/internal/miniostorage"
	"github.com/rs/zerolog"
	"io"
)

type fileService struct {
	logger        zerolog.Logger
	minioInstance *miniostorage.Instance
	connection    *pgx.Conn
}

func New(logger zerolog.Logger, minioInstance *miniostorage.Instance, connection *pgx.Conn) *fileService {
	return &fileService{logger: logger, minioInstance: minioInstance, connection: connection}
}

func (f *fileService) Store(entry *distrybute.FileEntry, reader io.Reader) (err error) {
	panic("implement me")
}

func (f *fileService) Request(callReference string) (entry distrybute.FileEntry, err error) {
	panic("implement me")
}

func (f *fileService) Delete(entries chan string) (responseChan chan distrybute.DeleteResponse) {
	panic("implement me")
}

func (f *fileService) ListEntries(limit int, offset int, sortBy distrybute.FileEntrySortElem, sortOrder distrybute.SortSequence, uuids ...uuid.UUID) (entries []distrybute.FileEntry, err error) {
	panic("implement me")
}

func (f *fileService) SearchEntries(query string, limit int, offset int, sortBy distrybute.FileEntrySortElem, sortOrder distrybute.SortSequence, uuids ...uuid.UUID) (entries []distrybute.FileEntry, err error) {
	panic("implement me")
}

func (f *fileService) ResolveMIMETypeStatistic(uuids ...uuid.UUID) (totalEntries int64, statistic distrybute.MIMETypeStatistic, err error) {
	panic("implement me")
}

func (f *fileService) ResolveUserUploadPeriodStatistic(uuid uuid.UUID, period distrybute.Period) (statistic distrybute.UserUploadPeriodStatistic, err error) {
	panic("implement me")
}
