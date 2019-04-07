package gosharexserver

import (
	"github.com/google/uuid"
	"time"
)

// DefaultUserId is used for old file entries with no specified author.
var DefaultUserId, _ = uuid.FromBytes(make([]byte, 16))

// PasswordHashAlgorithm is used to represent a password hashing algorithm in order to allow
// multiple different hashing implementations.
type PasswordHashAlgorithm string

const (
	// Unique strings for all password hashing algorithms.
	HashingArgon2ID PasswordHashAlgorithm = "argon2id"
	// LatestPasswordHashAlgorithm declares the default used and latest password hash algorithm.
	LatestPasswordHashAlgorithm PasswordHashAlgorithm = HashingArgon2ID
)

// Role wraps the type of roles
type Role string

const (
	// Role constants
	RoleAdmin Role = "ADMIN"
	RoleUser  Role = "USER"
)

// User contains the basic user data.
type User struct {
	// Id is a unqiue id which can be used to identify the user.
	Id uuid.UUID
	// Username is the name of the user to e.g. login with.
	Username string
	// Role indicates the status of the user inside the system and whether he has extended access to
	// other ressources.
	Role Role
	// AuthorizationToken holds the auth token for the user to use when uploading file entries.
	AuthorizationToken string
	// HashingAlgorithm indicates the hashing algorithm by providing an identical and unique number
	// to match the algorithm.
	PasswordHashAlgorithm PasswordHashAlgorithm
}

// IsUsingLatestPasswordHashAlgorithm indicates whether the user is using the latest password hash
// algorithm (LatestPasswordHashAlgorithm).
func (user User) IsUsingLatestPasswordHashAlgorithm() bool {
	if user.PasswordHashAlgorithm == "" {
		panic("could not check for password hash algorithm because no algorithm set in user instance")
	}
	return user.PasswordHashAlgorithm == LatestPasswordHashAlgorithm
}

// UserService contains the basic functions for interacting with the user database and their passwords.
type UserService interface {
	// CreateNewUser creates a new user by using the specified Username and role within the user
	// instance. After a successful creation the Id and AuthorizationToken of the user instance are
	// updated. It returns an error (err) if something went wrong.
	CreateNewUser(user *User, password []byte) (err error)
	// ResolveUser tries to search for the user by using the uuid or username set within the user instance.
	// After successfully finding the entry it sets the Id, Username, Role and PasswordHashAlgorithm
	// field of the user. It returns an error (err) if something went wrong.
	ResolveUser(user *User) (err error)
	// UpdateUsername updates the user`s username and sets the value of the user instance. It
	// returns an error (err) if something went wrong.
	UpdateUsername(user *User, newUsername string) (err error)
	// UpdateRole updates the user`s role and sets the value of the user instance. It returns an
	// error (err) if something went wrong.
	UpdateRole(user *User, newRole Role) (err error)
	// ResolveAuthorizationToken resolves the authorization token and sets the value of the user
	// instance. It returns an error (err) if something went wrong.
	ResolveAuthorizationToken(user *User) (err error)
	// UpdateAuthorizationToken updates the user`s authorization token and sets the value of the
	// user instance. It returns an error (err) if something went wrong.
	UpdateAuthorizationToken(user *User) (err error)
	// DeleteUser deletes the user by searching for the user`s id. It returns an error (err) if
	// something went wrong.
	DeleteUser(id uuid.UUID) (err error)
	// CheckPassword checks the user`s password and whether the username is existent inside the
	// database. It returns the User instance (user) with all values filled except for both password
	// fields if the check was successful. It returns an error (err) if something went wrong.
	CheckPassword(user *User, password []byte) (ok bool, err error)
	// UpdatePasssword updates the user`s password. It returns an error (err) if something went wrong.
	UpdatePasssword(user *User, password []byte) (err error)
}

var (
	// SessionExpiryPermanent is the value for the ExpiresOn field which represents a permanent session.
	SessionExpiryPermanent = time.Unix(-1, 0)
	// SessionExpiryOnyTimeOnly is the value for the ExpiresOn field which represents a one-time-only
	// session.
	SessionExpiryOnyTimeOnly = time.Unix(0, 0)
)

// SessionService contains the basic functions to manage sessions and serialize/deserialize them.
type SessionService interface {
	// NewUserSession returns a new session and saves it to the database. It returns an error if something
	NewUserSession(user *User, expiresOn time.Time) (session *UserSession, err error)
}

// UserSession contains the basic information needed within a session object. It does not contain any
// serialization functions - these belong to the SessionService implementation.
type UserSession struct {
	// ExpiresOn holds the date when the session should expire. If ExpiresOn#Unix is equals to -1,
	// the session should be permanent and if its equals to 0, it should be one session only.
	ExpiresOn time.Time
	// UserId holds the user`s id inside the database.
	UserId uuid.UUID
	// Username holds the user`s username.
	Username string
	// Role holds the user`s role.
	Role Role
	// Valid indicates whether the session is stilled marked as valid. This field is only used if
	// the session should be invalidated manually.
	Valid bool
}

// GetUserInstance returns a new instantiated instance of the User type by using the given data.
func (session UserSession) GetUserInstance() (user *User) {
	return &User{
		Id:       session.UserId,
		Username: session.Username,
		Role:     session.Role,
	}
}

// IsActive indicates whether the session is still active and marked as valid.
func (session UserSession) IsActiveAndValid() bool {
	return session.Valid && (session.ExpiresOn.Unix() == -1 || time.Now().Before(session.ExpiresOn))
}
