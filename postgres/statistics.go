package postgres

import (
	"github.com/google/uuid"
	"github.com/mmichaelb/gosharexserver"
	"time"
)

// ResolveMIMETypeStatistic is the Postgres based implementation of the Service interface ResolveMIMETypeStatistic function.
func (service *Service) ResolveMIMETypeStatistic(uid []uuid.UUID) (totalEntries int64, statistic gosharexserver.MIMETypeStatistic, err error) {
	panic("not implemented")
}

// UploadStatistics is the Postgres based implementation of the Service interface UploadStatistics function.
func (service *Service) UploadStatistics(uid []uuid.UUID, period time.Duration) (statistic gosharexserver.UserUploadPeriodStatistic, err error) {
	panic("not implemented")
}
