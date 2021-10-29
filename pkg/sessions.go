package pkg

import (
	"net/http"
)

// SessionService contains the basic functions to manage sessions.
type SessionService interface {
	// SetUserSession sets the session and saves it to the database. It returns an error if
	// something went wrong.
	SetUserSession(user *User, writer http.ResponseWriter) error
	// InvalidateUserSessions invalidates all used user sessions and therefore automatically logs
	// the user out of his account.
	InvalidateUserSessions(user *User) error
	// ValidateUserSession validates the http request and checks whether a session is present. If so,
	// the matched user is returned.
	ValidateUserSession(req *http.Request) (bool, *User, error)
}
