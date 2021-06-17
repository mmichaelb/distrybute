package postgres

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog/log"
)

func (s *service) initUserDDL() (err error) {
	rows, err := s.connection.Query(context.Background(), `CREATE TABLE distrybute.users (
		id uuid,
		username varchar(16) NOT NULL,
		auth_token text NULL,
		password_alg varchar(32) NOT NULL,
		password_salt bytea NOT NULL,
		"password" bytea NOT NULL,
		CONSTRAINT users_pk PRIMARY KEY (id),
		CONSTRAINT users_auth_token_unique UNIQUE (auth_token),
		CONSTRAINT users_username_un UNIQUE (username)
	)`)
	if err != nil {
		log.Err(err).Msg("could not run initial user ddl")
		return err
	}
	rows.Close()
	return nil
}

func (s *service) CreateNewUser(username string, password []byte) (user *distrybute.User, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	passwordAlgorithm := distrybute.LatestPasswordHashAlgorithm
	hashedPassword, salt, err := generatePasswordUserEntry(password, passwordAlgorithm)
	if err != nil {
		return nil, err
	}
	authToken, err := generateAuthToken()
	if err != nil {
		return nil, err
	}
	rows, err := s.connection.Query(context.Background(),
		`INSERT INTO distrybute.users (id, username, auth_token, password_alg, password_salt, password) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, username, authToken, string(passwordAlgorithm), salt, hashedPassword)
	if err != nil {
		return nil, err
	}
	rows.Close()
	return &distrybute.User{
		ID:                    id,
		Username:              username,
		AuthorizationToken:    authToken,
		PasswordHashAlgorithm: passwordAlgorithm,
	}, nil
}

func (s *service) CheckPassword(username string, password []byte) (ok bool, user *distrybute.User, err error) {
	row := s.connection.QueryRow(context.Background(),
		`SELECT (id, username, password, password_alg, password_salt) FROM distrybute.users WHERE user LIKE $1`, username)
	var id uuid.UUID
	var fetchedUsername string
	var expectedPasswordHash, passwordSalt []byte
	var passwordAlgorithm distrybute.PasswordHashAlgorithm
	err = row.Scan(&id, &fetchedUsername, &expectedPasswordHash, &passwordAlgorithm, &passwordSalt)
	if err == pgx.ErrNoRows {
		return false, nil, distrybute.ErrUserNotFound
	} else if err != nil {
		return false, nil, err
	}
	hashedPassword, err := generatePasswordHash(password, passwordSalt, passwordAlgorithm)
	if err != nil {
		return false, nil, err
	}
	if !bytes.Equal(expectedPasswordHash, hashedPassword) {
		return false, nil, nil
	}
	return true, &distrybute.User{
		ID:                    id,
		Username:              username,
		PasswordHashAlgorithm: passwordAlgorithm,
	}, nil
}

func (s *service) UpdateUsername(id uuid.UUID, newUsername string) (err error) {
	row := s.connection.QueryRow(context.Background(), `UPDATE distrybute.users SET username=$1 WHERE id=$2`, newUsername, id)
	err = row.Scan()
	if err == nil {
		return
	}
	pgErr, ok := err.(*pgconn.PgError)
	// check for unique constraint violation
	if ok && pgErr.Code == "23505" {
		return distrybute.ErrUserAlreadyExists
	} else {
		return err
	}
}

func (s *service) ResolveAuthorizationToken(id uuid.UUID) (token string, err error) {
	panic("implement me")
}

func (s *service) RefreshAuthorizationToken(id uuid.UUID) (token string, err error) {
	panic("implement me")
}

func (s *service) DeleteUser(id uuid.UUID) (err error) {
	panic("implement me")
}

func (s *service) UpdatePassword(id uuid.UUID, password []byte) (err error) {
	panic("implement me")
}
