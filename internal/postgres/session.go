package postgres

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type userKey struct{}

var sessionDuration = time.Hour * 24 * 7

const sessionKeyLength = 32

func (s *service) initSessionDDL() (err error) {
	row := s.connection.QueryRow(context.Background(), `CREATE TABLE IF NOT EXISTS distrybute.sessions (
		id uuid,
		key varchar(32),
		created_at timestamptz,
		valid_until timestamptz
	)`)
	if err = row.Scan(); err != nil {
		log.Err(err).Msg("could not run initial session ddl")
		return err
	}
	return nil
}

func (s *service) SetUserSession(user *distrybute.User, resp http.ResponseWriter) (err error) {

}

func (s *service) InvalidateUserSessions(user *distrybute.User) (err error) {
	panic("implement me")
}

func (s *service) ValidateUserSession(req *http.Request) (user *distrybute.User, err error) {
	panic("implement me")
}

func setUserContextValue() {

}

func (s *service) GetUserFromContext(req *http.Request) (user *distrybute.User) {
	panic("implement me")
}

func generateSessionKey() (key string, err error) {
	keyBytes := make([]byte, 16)
	_, err = rand.Read(keyBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}
