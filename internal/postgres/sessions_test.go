package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_generateSessionKey(t *testing.T) {
	t.Run("generates a session with a correct length", func(t *testing.T) {
		key, err := generateSessionKey()
		assert.NoError(t, err)
		assert.Len(t, key, sessionKeyLength)
	})
	t.Run("generates different session keys", func(t *testing.T) {
		gotFirstKey, err := generateSessionKey()
		assert.NoError(t, err)
		gotSecondKey, err := generateSessionKey()
		assert.NoError(t, err)
		assert.True(t, gotFirstKey != gotSecondKey)
	})
}
