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
			const username = "usertest-duplicate"
			_, err := userService.CreateNewUser(username, []byte("Sommer2019"))
			assert.NoError(t, err)
			_, err = userService.CreateNewUser(username, []byte("Sommer2020"))
			assert.ErrorIs(t, err, distrybute.ErrUserAlreadyExists)
		})
		t.Run("duplicate usernames are being detected case insensitively", func(t *testing.T) {
			const username = "usertest-duplicate-case-ins"
			const usernameSecond = "usertest-duplicate-case-INS"
			_, err := userService.CreateNewUser(username, []byte("Sommer2019"))
			assert.NoError(t, err)
			_, err = userService.CreateNewUser(usernameSecond, []byte("Sommer2020"))
			assert.ErrorIs(t, err, distrybute.ErrUserAlreadyExists)
		})
		t.Run("user is deleted correctly", func(t *testing.T) {
			const username = "usertest-del"
			user, err := userService.CreateNewUser(username, []byte("Testpassword"))
			assert.NoError(t, err)
			err = userService.DeleteUser(user.ID)
			assert.NoError(t, err)
			_, _, err = userService.CheckPassword(username, []byte("Testpassword"))
			assert.ErrorIs(t, err, distrybute.ErrUserNotFound)
		})
		t.Run("user cannot be deleted if not present", func(t *testing.T) {
			id, err := uuid.Parse("7c478fdc-be22-4571-b7b6-2dfa5a31a1a7") // parse some random uuid
			assert.NotNil(t, err, "uuid could not be parsed")
			err = userService.DeleteUser(id)
			assert.ErrorIs(t, err, distrybute.ErrUserNotFound)
		})
	}
}
