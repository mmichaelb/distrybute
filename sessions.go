package gosharexserver

import (
	"errors"
	"net/http"
)

var (
	// ErrSessionInvalid indicates that the given session has been invalidated.
	ErrSessionInvalid = errors.New("the given session is invalid")
	// ErrNoSessionSet indicates that the given request does not contain a session.
	ErrNoSessionSet = errors.New("the given request does not contain a session")
)

// SessionService contains the basic functions to manage sessions.
type SessionService interface {
	// SetUserSession sets the session and saves it to the database. It returns an error if
	// something went wrong.
	SetUserSession(user *User, resp http.ResponseWriter) (err error)
	// InvalidateUserSessions invalidates all used user sessions and therefore automatically logs
	// the user out of his account.
	InvalidateUserSessions(user *User) (err error)
	// ValidateUserSession validates the http request and checks whether a session is present. If so,
	// the matched user is returned.
	ValidateUserSession(req *http.Request) (user *User, err error)
}
