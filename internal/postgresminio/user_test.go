package postgresminio

import (
	"github.com/google/uuid"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/stretchr/testify/assert"
	"testing"
)

func userServiceIntegrationTest(userService distrybute.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("user is created correctly", func(t *testing.T) {
			const username = "usertest-user"
			user, err := userService.CreateNewUser(username, []byte("somepassword"))
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.UUID{}, user.ID)
			assert.Equal(t, username, user.Username)
			assert.NotEmpty(t, user.AuthorizationToken)
			assert.Equal(t, distrybute.LatestPasswordHashAlgorithm, user.PasswordHashAlgorithm)
		})
		t.Run("duplicate usernames are not accepted", func(t *testing.T) {
			const username = "usertest-dupl"
			_, err := userService.CreateNewUser(username, []byte("Sommer2019"))
			assert.NoError(t, err)
			_, err = userService.CreateNewUser(username, []byte("Sommer2020"))
			assert.ErrorIs(t, err, distrybute.ErrUserAlreadyExists)
		})
	}
}
