package postgres

import (
	"encoding/hex"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generateAuthToken(t *testing.T) {
	t.Run("returns an encoded hex value", func(t *testing.T) {
		gotAuthToken, err := generateAuthToken()
		assert.NoError(t, err)
		_, decodeErr := hex.DecodeString(gotAuthToken)
		assert.NoError(t, decodeErr)
	})
	t.Run("returns a long enough hex value", func(t *testing.T) {
		gotAuthToken, err := generateAuthToken()
		assert.NoError(t, err)
		tokenBytes, decodeErr := hex.DecodeString(gotAuthToken)
		assert.NoError(t, decodeErr)
		assert.Len(t, tokenBytes, 16)
	})
	t.Run("generates different auth tokens", func(t *testing.T) {
		firstGotAuthToken, err := generateAuthToken()
		assert.NoError(t, err)
		secondGotAuthToken, err := generateAuthToken()
		assert.NoError(t, err)
		assert.False(t, firstGotAuthToken == secondGotAuthToken)
	})
}

func Test_generatePasswordHash(t *testing.T) {
	t.Run("argon2id hash being generated correctly", func(t *testing.T) {
		gotPasswordHash, err := generatePasswordHash([]byte("Sommer2019"), []byte("somegoodsalt"), distrybute.PasswordHashArgon2ID)
		assert.NoError(t, err)
		expectedHash := []byte{0xf5, 0x5f, 0xae, 0xf, 0xbd, 0x24, 0x81, 0x8e, 0xe5, 0xb7, 0x14, 0x7e, 0xee, 0x98, 0xa6, 0x50, 0xc3, 0xbc, 0xd1, 0x3, 0x34, 0xcb, 0xc8, 0x2b, 0x29, 0x44, 0x9c, 0x64, 0x2d, 0x22, 0xa8, 0x9d}
		assert.Equal(t, expectedHash, gotPasswordHash)
	})
	t.Run("unknown password hash algorithm is being detected", func(t *testing.T) {
		_, err := generatePasswordHash([]byte("Sommer2019"), []byte("somegoodsalt"), "notapasswordhashalgorithmatall")
		assert.Error(t, err)
	})
}

func Test_generatePasswordUserEntry(t *testing.T) {
	t.Run("salt being generated is of correct length", func(t *testing.T) {
		_, salt, err := generatePasswordUserEntry([]byte("Sommer2019"), distrybute.PasswordHashArgon2ID)
		assert.NoError(t, err)
		assert.Len(t, salt, 16)
	})
	t.Run("password hash is being generated correctly", func(t *testing.T) {
		password := []byte("Sommer2019")
		gotPasswordHash, salt, err := generatePasswordUserEntry(password, distrybute.PasswordHashArgon2ID)
		assert.NoError(t, err)
		expectedPasswordHash, err := generatePasswordHash(password, salt, distrybute.PasswordHashArgon2ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedPasswordHash, gotPasswordHash)
	})
}
