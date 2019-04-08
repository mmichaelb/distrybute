package gosharexserver

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
)

var (
	// ErrSessionInvalid indicates that the given session has been invalidated.
	ErrSessionInvalid = errors.New("the given session is invalid")
)

// SessionService contains the basic functions to manage sessions and serialize/deserialize them.
type SessionService interface {
	// NewUserSession returns a new session and saves it to the database. It returns an error if
	// something went wrong.
	NewUserSession(user *User, expiresOn time.Time) (session *UserSession, err error)
	// InvalidateUserSessions invalidates all used user sessions and therefore automatically logs
	// the user out of his account.
	InvalidateUserSessions(user *User) (err error)
	// ParseSessionFromRequest parses the user session from the given http request and returns an
	// error if something went wrong. If the error is an ErrSessionInvalid, the session will still
	// be parsed but it is not valid anymore.
	ParseSessionFromRequest(req *http.Request) (session *UserSession, err error)
}

// UserSession contains the basic information needed within a session object. It does not contain any
// serialization functions - these belong to the SessionService implementation.
type UserSession struct {
	// ExpiresOn holds the date when the session should expire.
	ExpiresOn time.Time
	// UserId holds the user`s id inside the database.
	UserID uuid.UUID
	// Username holds the user`s username.
	Username string
	// Role holds the user`s role.
	Role Role
}

// GetUserInstance returns a new instantiated instance of the User type by using the given data.
func (session UserSession) GetUserInstance() (user *User) {
	return &User{
		Id:       session.UserID,
		Username: session.Username,
		Role:     session.Role,
	}
}

// IsActive returns whether the session is still active.
func (session UserSession) IsActive() (active bool) {
	return time.Now().Before(session.ExpiresOn)
}
