package postgresminio

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mmichaelb/distrybute/pkg"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func userServiceIntegrationTest(userService distrybute.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		const usernameBasicTest = "usertest-user"
		var userBasicTestId uuid.UUID
		t.Run("user is created correctly", func(t *testing.T) {
			user, err := userService.CreateNewUser(usernameBasicTest, []byte("somepassword"))
			assert.NoError(t, err)
			assert.NotEqual(t, uuid.UUID{}, user.ID)
			userBasicTestId = user.ID
			assert.Equal(t, usernameBasicTest, user.Username)
			assert.NotEmpty(t, user.AuthorizationToken)
			assert.Equal(t, distrybute.LatestPasswordHashAlgorithm, user.PasswordHashAlgorithm)
		})
		t.Run("duplicate usernames are not accepted", func(t *testing.T) {
			_, err := userService.CreateNewUser(usernameBasicTest, []byte("Sommer2020"))
			assert.ErrorIs(t, err, distrybute.ErrUserAlreadyExists)
		})
		t.Run("duplicate usernames are being detected case insensitively", func(t *testing.T) {
			_, err := userService.CreateNewUser(strings.ToUpper(usernameBasicTest), []byte("Sommer2020"))
			assert.ErrorIs(t, err, distrybute.ErrUserAlreadyExists)
		})
		t.Run("user deletion test", func(t *testing.T) {
			t.Run("user is deleted correctly", func(t *testing.T) {
				err := userService.DeleteUser(userBasicTestId)
				assert.NoError(t, err)
				_, _, err = userService.CheckPassword(usernameBasicTest, []byte("Testpassword"))
				assert.ErrorIs(t, err, distrybute.ErrUserNotFound)
			})
			t.Run("user cannot be deleted if not present", func(t *testing.T) {
				id, err := uuid.Parse("7c478fdc-be22-4571-b7b6-2dfa5a31a1a7") // parse some random uuid
				assert.Nil(t, err, "uuid could not be parsed")
				err = userService.DeleteUser(id)
				assert.ErrorIs(t, err, distrybute.ErrUserNotFound)
			})
		})
		t.Run("user list tests", func(t *testing.T) {
			const userAmount = 5
			const usernamePattern = "usertest-list-%d"
			users := make([]*distrybute.User, userAmount)
			for i := 0; i < userAmount; i++ {
				user, err := userService.CreateNewUser(fmt.Sprintf(usernamePattern, i), []byte("testpassowrd"))
				assert.NoError(t, err)
				users[i] = user
				t.Cleanup(func() {
					_ = userService.DeleteUser(user.ID)
				})
			}
			retrievedUsers, err := userService.ListUsers()
			assert.NoError(t, err, "list users method returned a non-nil err")
			assert.Len(t, retrievedUsers, userAmount)
			for _, createdUser := range users {
				found := false
				for _, retrievedUser := range retrievedUsers {
					if createdUser.ID == retrievedUser.ID {
						found = true
					}
				}
				assert.True(t, found, fmt.Sprintf("user %s:%s could not be found within the returned user list (len: %d - %v)",
					createdUser.Username, createdUser.ID, len(retrievedUsers), retrievedUsers))
			}
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
				assert.ErrorIs(t, err, distrybute.ErrUserNotFound)
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
			assert.ErrorIs(t, err, distrybute.ErrUserNotFound, "login with old username was still successful")
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
		t.Run("user retrieved by username tests", func(t *testing.T) {
			const username = "usertest-retrieve-by-username"
			password := []byte("Sommer2019")
			user, err := userService.CreateNewUser(username, password)
			assert.NoError(t, err)
			t.Run("user can be retrieved by using the username", func(t *testing.T) {
				retrievedUser, err := userService.GetUserByUsername(username)
				assert.NoError(t, err, "user could not be resolved")
				assert.Equal(t, user.ID, retrievedUser.ID)
			})
			t.Run("user can be retrieved by using a username case insensitively", func(t *testing.T) {
				retrievedUser, err := userService.GetUserByUsername(strings.ToUpper(username))
				assert.NoError(t, err, "user could not be resolved case insensitively")
				assert.Equal(t, user.ID, retrievedUser.ID)
				assert.Equal(t, user.Username, retrievedUser.Username)
			})
			t.Run("no user can be found using a non-existent username", func(t *testing.T) {
				retrievedUser, err := userService.GetUserByUsername("this-user-does-not-exist")
				assert.ErrorIs(t, distrybute.ErrUserNotFound, err, "no error was returned when searching for non-existent user")
				assert.Nil(t, retrievedUser, "returned user is not nil")
			})
		})
	}
}
