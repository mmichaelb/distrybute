package gosharexserver

import (
	"github.com/google/uuid"
)

// DefaultUserId is used for old file entries with no specified author.
var DefaultUserId, _ = uuid.FromBytes(make([]byte, 16))
