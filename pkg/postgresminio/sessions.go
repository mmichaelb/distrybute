package postgresminio

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/mmichaelb/distrybute/pkg"
	"net/http"
	"time"
)

var sessionDuration = time.Hour * 24 * 7

const (
	sessionKeyLength  = 32
	sessionCookieName = "session_key"
)

func (s *service) initSessionDDL() (err error) {
	row := s.connection.QueryRow(context.Background(), `CREATE TABLE IF NOT EXISTS distrybute.sessions (
		id uuid NOT NULL,
		session_key varchar(32) NULL,
		created_at timestamptz NULL,
		valid_until timestamptz NULL,
		CONSTRAINT sessions_pk PRIMARY KEY (id),
		CONSTRAINT sessions_fk FOREIGN KEY (id) REFERENCES distrybute.users(id),
		CONSTRAINT sessions_session_key_unique UNIQUE (session_key)
	);`)
	if err = row.Scan(); !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("error occurred while running session ddl: %w", err)
	}
	return nil
}

func (s *service) SetUserSession(user *pkg.User, writer http.ResponseWriter) error {
	sessionKey, err := generateSessionKey()
	if err != nil {
		return err
	}
	createdAt := time.Now()
	validUntil := createdAt.Add(sessionDuration)
	row := s.connection.QueryRow(context.Background(),
		`INSERT INTO distrybute.sessions (id, session_key, created_at, valid_until) VALUES ($1, $2, $3, $4) 
			ON CONFLICT (id) DO UPDATE SET session_key=$2, created_at=$3, valid_until=$4`,
		user.ID, sessionKey, createdAt, validUntil)
	if err = row.Scan(); !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	cookie := &http.Cookie{
		Name:     "session_key",
		Value:    sessionKey,
		Path:     "/",
		Expires:  validUntil,
		Secure:   true,
		HttpOnly: true,
	}
	// set cookie
	http.SetCookie(writer, cookie)
	return nil
}

func (s *service) InvalidateUserSessions(user *pkg.User) (err error) {
	row := s.connection.QueryRow(context.Background(), `DELETE FROM distrybute.sessions WHERE id=$1`, user.ID)
	if err := row.Scan(); !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	return nil
}

func (s *service) ValidateUserSession(req *http.Request) (bool, *pkg.User, error) {
	cookie, err := req.Cookie(sessionCookieName)
	if errors.Is(err, http.ErrNoCookie) {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}
	sessionKey := cookie.Value
	row := s.connection.QueryRow(context.Background(),
		`SELECT id, username FROM distrybute.users WHERE id=(SELECT id FROM distrybute.sessions WHERE session_key LIKE $1)`,
		sessionKey)
	var id uuid.UUID
	var username string
	if err = row.Scan(&id, &username); errors.Is(err, pgx.ErrNoRows) {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}
	user := &pkg.User{
		ID:       id,
		Username: username,
	}
	return true, user, nil
}

func generateSessionKey() (key string, err error) {
	keyBytes := make([]byte, 16)
	_, err = rand.Read(keyBytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(keyBytes), nil
}
