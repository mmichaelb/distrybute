package postgres

import (
	"github.com/google/uuid"
	"github.com/mmichaelb/gosharexserver"
	"time"
)

// ResolveMIMETypeStatistic is the Postgres based implementation of the Manager interface ResolveMIMETypeStatistic function.
func (manager *Manager) ResolveMIMETypeStatistic(uid []uuid.UUID) (totalEntries int64, statistic gosharexserver.MIMETypeStatistic, err error) {
	panic("not implemented")
}

// UploadStatistics is the Postgres based implementation of the Manager interface UploadStatistics function.
func (manager *Manager) UploadStatistics(uid []uuid.UUID, period time.Duration) (statistic gosharexserver.UserUploadPeriodStatistic, err error) {
	panic("not implemented")
}
