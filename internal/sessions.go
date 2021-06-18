package distrybute

import (
	"net/http"
)

// SessionService contains the basic functions to manage sessions.
type SessionService interface {
	// SetUserSession sets the session and saves it to the database. It returns an error if
	// something went wrong. Moreover, this method also sets the request context to use when calling GetUserFromContext.
	SetUserSession(user *User, req *http.Request, writer http.ResponseWriter) (*http.Request, error)
	// InvalidateUserSessions invalidates all used user sessions and therefore automatically logs
	// the user out of his account.
	InvalidateUserSessions(user *User) (err error)
	// ValidateUserSession validates the http request and checks whether a session is present. If so,
	// the matched user is returned. It also sets the http request context value to use when calling GetUserFromContext.
	ValidateUserSession(req *http.Request) (bool, *http.Request, error)
	// GetUserFromContext returns the user bound to the http request context. Returns nil, if no user is bound.
	GetUserFromContext(req *http.Request) (user *User)
}
