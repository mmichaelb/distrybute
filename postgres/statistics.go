package postgres

import (
	"github.com/google/uuid"
	"github.com/mmichaelb/gosharexserver"
	"time"
)

func (manager *Manager) ContentTypeStatistics(uid []uuid.UUID) (err error, stats *gosharexserver.MapBasedStatistics) {
	panic("not implemented")
}

func (manager *Manager) UploadStatistics(uid []uuid.UUID, period time.Duration) (err error, stats *gosharexserver.MapBasedStatistics) {
	panic("not implemented")
}
