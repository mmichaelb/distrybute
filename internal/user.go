package distrybute

import (
	"errors"
	"github.com/google/uuid"
)

// DefaultUserID is used for old file entries with no specified author.
var DefaultUserID = uuid.UUID{}

// PasswordHashAlgorithm is used to represent a password hashing algorithm in order to allow
// multiple different hashing implementations.
type PasswordHashAlgorithm string

const (
	// Unique strings for all password hashing algorithms.

	// PasswordHashArgon2ID is the identical name for the expensive key derivation function Argon2Id.
	PasswordHashArgon2ID PasswordHashAlgorithm = "argon2id"
	// LatestPasswordHashAlgorithm declares the default used and latest password hash algorithm.
	LatestPasswordHashAlgorithm = PasswordHashArgon2ID
)

// User contains the basic user data.
type User struct {
	// ID is a unqiue ID which can be used to identify the user.
	ID uuid.UUID
	// Username is the name of the user to e.g. login with.
	Username string
	// AuthorizationToken holds the auth token for the user to use when uploading file entries.
	AuthorizationToken string
	// PasswordHashAlgorithm indicates the hashing algorithm which this user entry is using.
	PasswordHashAlgorithm PasswordHashAlgorithm
}

// IsUsingLatestPasswordHashAlgorithm indicates whether the user is using the latest password hash
// algorithm (LatestPasswordHashAlgorithm).
func (user User) IsUsingLatestPasswordHashAlgorithm() (bool, error) {
	if user.PasswordHashAlgorithm == "" {
		return false, errors.New("could not check for password hash algorithm because no algorithm set in user instance")
	}
	return user.PasswordHashAlgorithm == LatestPasswordHashAlgorithm, nil
}

var (
	ErrUserAlreadyExists = errors.New("the user already exists")
	ErrUserNotFound      = errors.New("the given user could not be found")
)

// UserService contains the basic functions for interacting with the user database and their passwords.
type UserService interface {
	// CreateNewUser creates a new user by using the specified Username. After a successful creation, a user instance
	// is returned. It returns an error (err) if something went wrong.
	CreateNewUser(username string, password []byte) (user *User, err error)
	// CheckPassword checks the user`s password and whether the username is existent inside the database. Ok is true if
	// the check was successful. If the user could not be found a ErrUserNotFound is returned.
	CheckPassword(username string, password []byte) (ok bool, err error)
	// UpdateUsername updates the user`s username and sets the value of the user instance. It
	// returns an error (err) if something went wrong.
	UpdateUsername(user *User, newUsername string) (err error)
	// ResolveAuthorizationToken resolves the authorization token and sets the value of the user
	// instance. It returns an error (err) if something went wrong.
	ResolveAuthorizationToken(user *User) (err error)
	// UpdateAuthorizationToken updates the user`s authorization token and sets the value of the
	// user instance. It returns an error (err) if something went wrong.
	UpdateAuthorizationToken(user *User) (err error)
	// DeleteUser deletes the user by searching for the user`s ID. It returns an error (err) if
	// something went wrong.
	DeleteUser(id uuid.UUID) (err error)
	// UpdatePassword updates the user`s password. It returns an error (err) if something went wrong.
	UpdatePassword(user *User, password []byte) (err error)
}
