package distrybute

import (
	"github.com/google/uuid"
)

// DefaultUserID is used for old file entries with no specified author.
var DefaultUserID, _ = uuid.FromBytes(make([]byte, 16))

// PasswordHashAlgorithm is used to represent a password hashing algorithm in order to allow
// multiple different hashing implementations.
type PasswordHashAlgorithm string

const (
	// Unique strings for all password hashing algorithms.

	// PasswordHashArgon2ID is the identical name for the expensive key derivation function Argon2Id.
	PasswordHashArgon2ID PasswordHashAlgorithm = "argon2id"
	// LatestPasswordHashAlgorithm declares the default used and latest password hash algorithm.
	LatestPasswordHashAlgorithm PasswordHashAlgorithm = PasswordHashArgon2ID
)

// Role is the string wrapper for user roles.
type Role string

const (
	// RoleAdmin represents the constant for the admin role.
	RoleAdmin Role = "ADMIN"
	// RoleUser represents the constant for the user role.
	RoleUser Role = "USER"
)

// User contains the basic user data.
type User struct {
	// ID is a unqiue ID which can be used to identify the user.
	ID uuid.UUID
	// Username is the name of the user to e.g. login with.
	Username string
	// Role indicates the status of the user inside the system and whether he has extended access to
	// other ressources.
	Role Role
	// AuthorizationToken holds the auth token for the user to use when uploading file entries.
	AuthorizationToken string
	// PasswordHashAlgorithm indicates the hashing algorithm which this user entry is using.
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
	// DeleteUser deletes the user by searching for the user`s ID. It returns an error (err) if
	// something went wrong.
	DeleteUser(id uuid.UUID) (err error)
	// CheckPassword checks the user`s password and whether the username is existent inside the
	// database. It returns the User instance (user) with all values filled except for both, the
	// authorization token and  password field if the check was successful. It returns an error
	// (err) if something went wrong.
	CheckPassword(user *User, password []byte) (ok bool, err error)
	// UpdatePassword updates the user`s password. It returns an error (err) if something went wrong.
	UpdatePassword(user *User, password []byte) (err error)
}
