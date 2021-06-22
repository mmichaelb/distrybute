package postgresminio

import (
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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

func sessionServiceIntegrationTest(sessionService distrybute.SessionService, userService distrybute.UserService) func(t *testing.T) {
	return func(t *testing.T) {
		user, err := userService.CreateNewUser("sessiontest-user", []byte("Sommer2019"))
		assert.NoError(t, err, "could not create session test user")
		t.Run("user session is set correctly", func(t *testing.T) {
			req := httptest.NewRequest("/", http.MethodGet, nil)
			writer := httptest.NewRecorder()
			req, err := sessionService.SetUserSession(user, req, writer)
			assert.NoError(t, err, "could not set user session")
			assert.Len(t, writer.Result().Cookies(), 1, "session key cookie is not set")
			cookie := writer.Result().Cookies()[0]
			assert.Equal(t, "session_key", cookie.Name, "session key is not named correctly")
			assert.True(t, cookie.HttpOnly, "cookie should be http only")
			assert.True(t, cookie.Secure, "cookie be secure (only sent via ssl)")
			assert.NotEmpty(t, cookie.Value, "cookie value should not be empty")
			assert.Equal(t, user, sessionService.GetUserFromContext(req), "user could not be retrieved from request")
		})
		t.Run("user session is parsed correctly", func(t *testing.T) {
			req := httptest.NewRequest("/", http.MethodGet, nil)
			writer := httptest.NewRecorder()
			_, err := sessionService.SetUserSession(user, req, writer)
			assert.NoError(t, err, "could not set user session")
			cookie := writer.Result().Cookies()[0]
			req = httptest.NewRequest("/", http.MethodGet, nil)
			req.AddCookie(cookie)
			ok, req, err := sessionService.ValidateUserSession(req)
			assert.True(t, ok, "session is detected as invalid")
			assert.NoError(t, err, "session could not be validated")
			contextUser := sessionService.GetUserFromContext(req)
			assert.Equal(t, user.ID, contextUser.ID)
			assert.Equal(t, user.Username, contextUser.Username)
		})
		t.Run("user session can be invalidated", func(t *testing.T) {
			req := httptest.NewRequest("/", http.MethodGet, nil)
			writer := httptest.NewRecorder()
			_, err := sessionService.SetUserSession(user, req, writer)
			assert.NoError(t, err, "could not set user session")
			cookie := writer.Result().Cookies()[0]
			err = sessionService.InvalidateUserSessions(user)
			req = httptest.NewRequest("/", http.MethodGet, nil)
			assert.NoError(t, err, "could not invalidate user session")
			req.AddCookie(cookie)
			ok, req, err := sessionService.ValidateUserSession(req)
			assert.False(t, ok, "session is still valid")
			assert.NoError(t, err, "could not validate session")
			user := sessionService.GetUserFromContext(req)
			assert.Nil(t, user, "user value is still set to context")
		})
	}
}
