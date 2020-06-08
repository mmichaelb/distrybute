package distrybute

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUser_IsUsingLatestPasswordHashAlgorithm(t *testing.T) {
	t.Run("empty password alg", func(t *testing.T) {
		user := &User{
			PasswordHashAlgorithm: "",
		}
		res, err := user.IsUsingLatestPasswordHashAlgorithm()
		assert.False(t, res)
		assert.Error(t, err)
	})
	t.Run("not latest password alg", func(t *testing.T) {
		user := &User{
			PasswordHashAlgorithm: "md5",
		}
		res, err := user.IsUsingLatestPasswordHashAlgorithm()
		assert.False(t, res)
		assert.NoError(t, err)
	})
	t.Run("latest password alg", func(t *testing.T) {
		user := &User{
			PasswordHashAlgorithm: LatestPasswordHashAlgorithm,
		}
		res, err := user.IsUsingLatestPasswordHashAlgorithm()
		assert.True(t, res)
		assert.NoError(t, err)
	})
}
