package postgresminio

import (
	"github.com/google/uuid"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func userServiceIntegrationTest(userService pkg.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		t.Run("user is created correctly", func(t *testing.T) {
			const username = "usertest-user"
			user, err := userService.CreateNewUser(username, []byte("somepassword"))
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.UUID{}, user.ID)
			assert.Equal(t, username, user.Username)
			assert.NotEmpty(t, user.AuthorizationToken)
			assert.Equal(t, pkg.LatestPasswordHashAlgorithm, user.PasswordHashAlgorithm)
		})
		t.Run("duplicate usernames are not accepted", func(t *testing.T) {
			const username = "usertest-duplicate"
			_, err := userService.CreateNewUser(username, []byte("Sommer2019"))
			assert.NoError(t, err)
			_, err = userService.CreateNewUser(username, []byte("Sommer2020"))
			assert.ErrorIs(t, err, pkg.ErrUserAlreadyExists)
		})
		t.Run("duplicate usernames are being detected case insensitively", func(t *testing.T) {
			const username = "usertest-duplicate-case-ins"
			const usernameSecond = "usertest-duplicate-case-INS"
			_, err := userService.CreateNewUser(username, []byte("Sommer2019"))
			assert.NoError(t, err)
			_, err = userService.CreateNewUser(usernameSecond, []byte("Sommer2020"))
			assert.ErrorIs(t, err, pkg.ErrUserAlreadyExists)
		})
		t.Run("user deletion test", func(t *testing.T) {
			t.Run("user is deleted correctly", func(t *testing.T) {
				const username = "usertest-del"
				user, err := userService.CreateNewUser(username, []byte("Testpassword"))
				assert.NoError(t, err)
				err = userService.DeleteUser(user.ID)
				assert.NoError(t, err)
				_, _, err = userService.CheckPassword(username, []byte("Testpassword"))
				assert.ErrorIs(t, err, pkg.ErrUserNotFound)
			})
			t.Run("user cannot be deleted if not present", func(t *testing.T) {
				id, err := uuid.Parse("7c478fdc-be22-4571-b7b6-2dfa5a31a1a7") // parse some random uuid
				assert.Nil(t, err, "uuid could not be parsed")
				err = userService.DeleteUser(id)
				assert.ErrorIs(t, err, pkg.ErrUserNotFound)
			})
		})
		t.Run("password check test", func(t *testing.T) {
			const username = "usertest-password-check"
			password := []byte("Sommer2019")
			user, err := userService.CreateNewUser(username, password)
			assert.NoError(t, err)
			t.Run("password check is done correctly", func(t *testing.T) {
				ok, resolvedUser, err := userService.CheckPassword(user.Username, password)
				assert.NoError(t, err)
				assert.True(t, ok)
				assert.Equal(t, user.ID, resolvedUser.ID)
				assert.Equal(t, user.Username, resolvedUser.Username)
				assert.Empty(t, resolvedUser.AuthorizationToken)
			})
			t.Run("password is checked correctly even if username is not of correct case", func(t *testing.T) {
				upperUsername := strings.ToUpper(user.Username)
				ok, resolvedUser, err := userService.CheckPassword(upperUsername, password)
				assert.NoError(t, err)
				assert.True(t, ok)
				assert.Equal(t, user.ID, resolvedUser.ID)
				assert.Equal(t, user.Username, resolvedUser.Username)
			})
			t.Run("wrong password is not accepted", func(t *testing.T) {
				ok, resolvedUser, err := userService.CheckPassword(username, []byte("nottherightpassword"))
				assert.NoError(t, err)
				assert.False(t, ok)
				assert.Nil(t, resolvedUser)
			})
			t.Run("username has to be registered within the system", func(t *testing.T) {
				ok, resolvedUser, err := userService.CheckPassword("userthatdoesnotexist", []byte("nottherightpassword"))
				assert.ErrorIs(t, err, pkg.ErrUserNotFound)
				assert.False(t, ok)
				assert.Nil(t, resolvedUser)
			})
		})
		t.Run("username is being updated correctly", func(t *testing.T) {
			const username = "usertest-update-username"
			password := []byte("Sommer2019")
			user, err := userService.CreateNewUser(username, password)
			assert.NoError(t, err)
			newUsername := "usertest-update-username-new"
			err = userService.UpdateUsername(user.ID, newUsername)
			assert.NoError(t, err, "username could not be updated")
			ok, resolvedUser, err := userService.CheckPassword(newUsername, password)
			assert.NoError(t, err, "could not check password with new username")
			assert.True(t, ok)
			assert.Equal(t, user.ID, resolvedUser.ID)
			assert.Equal(t, newUsername, resolvedUser.Username)
			ok, resolvedUser, err = userService.CheckPassword(username, password)
			assert.ErrorIs(t, err, pkg.ErrUserNotFound, "login with old username was still successful")
			assert.False(t, ok)
			assert.Nil(t, resolvedUser)
		})
		t.Run("authorization token tests", func(t *testing.T) {
			const username = "usertest-auth-token"
			password := []byte("Sommer2019")
			user, err := userService.CreateNewUser(username, password)
			assert.NoError(t, err)
			t.Run("authorization token can be retrieved", func(t *testing.T) {
				token, err := userService.ResolveAuthorizationToken(user.ID)
				assert.NoError(t, err, "authorization token could not be resolved")
				assert.Equal(t, user.AuthorizationToken, token)
			})
			t.Run("authorization token can be used to retrieve a user", func(t *testing.T) {
				ok, retrievedUser, err := userService.GetUserByAuthorizationToken(user.AuthorizationToken)
				assert.NoError(t, err, "authorization token could not be used to retrieve a user")
				assert.True(t, ok)
				assert.Equal(t, user.Username, retrievedUser.Username)
				assert.Equal(t, user.ID, retrievedUser.ID)
			})
			t.Run("authorization token can be refreshed", func(t *testing.T) {
				token, err := userService.RefreshAuthorizationToken(user.ID)
				assert.NoError(t, err, "authorization token could not be refreshed")
				assert.NotEqual(t, user.AuthorizationToken, token)
				retrievedToken, err := userService.ResolveAuthorizationToken(user.ID)
				assert.NoError(t, err, "authorization token could not be resolved")
				assert.Equal(t, retrievedToken, token)
			})
		})
	}
}
