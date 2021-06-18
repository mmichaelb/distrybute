package postgresminio

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"golang.org/x/crypto/argon2"
)

const (
	saltLength      = 16
	authTokenLength = 16
)

func generatePasswordUserEntry(password []byte, algorithm distrybute.PasswordHashAlgorithm) (hashedPassword []byte, salt []byte, err error) {
	salt = make([]byte, saltLength)
	if _, err = rand.Read(salt); err != nil {
		return
	}
	hashedPassword, err = generatePasswordHash(password, salt, algorithm)
	return hashedPassword, salt, err
}

func generatePasswordHash(password []byte, salt []byte, algorithm distrybute.PasswordHashAlgorithm) (hashedPassword []byte, err error) {
	switch algorithm {
	case distrybute.PasswordHashArgon2ID:
		hashedPassword = argon2.IDKey(password, salt, 1, 64*1024, 4, 32)
	default:
		return nil, errors.New("the provided hashing algorithm is unknown")
	}
	return
}

func generateAuthToken() (authToken string, err error) {
	authTokenBytes := make([]byte, authTokenLength)
	_, err = rand.Read(authTokenBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(authTokenBytes), nil
}
