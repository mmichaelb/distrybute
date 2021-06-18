package postgresminio

import (
	"github.com/google/uuid"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"io"
)

func (s *service) Store(entry *distrybute.FileEntry, reader io.Reader) (err error) {

}

func (s *service) Request(callReference string) (entry distrybute.FileEntry, err error) {
	panic("implement me")
}

func (s *service) Delete(entries chan string) (responseChan chan distrybute.DeleteResponse) {
	panic("implement me")
}

func (s *service) ListEntries(limit int, offset int, sortBy distrybute.FileEntrySortElem, sortOrder distrybute.SortSequence, uuids ...uuid.UUID) (entries []distrybute.FileEntry, err error) {
	panic("implement me")
}

func (s *service) SearchEntries(query string, limit int, offset int, sortBy distrybute.FileEntrySortElem, sortOrder distrybute.SortSequence, uuids ...uuid.UUID) (entries []distrybute.FileEntry, err error) {
	panic("implement me")
}

func (s *service) ResolveMIMETypeStatistic(uuids ...uuid.UUID) (totalEntries int64, statistic distrybute.MIMETypeStatistic, err error) {
	panic("implement me")
}

func (s *service) ResolveUserUploadPeriodStatistic(uuid uuid.UUID, period distrybute.Period) (statistic distrybute.UserUploadPeriodStatistic, err error) {
	panic("implement me")
}
