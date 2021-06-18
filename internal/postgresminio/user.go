package postgresminio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	distrybute "github.com/mmichaelb/distrybute/internal"
)

func (s *service) initUserDDL() (err error) {
	row := s.connection.QueryRow(context.Background(), `CREATE TABLE IF NOT EXISTS distrybute.users (
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
	if err = row.Scan(); !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("error occurred while running user ddl: %w", err)
	}
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
	row := s.connection.QueryRow(context.Background(),
		`INSERT INTO distrybute.users (id, username, auth_token, password_alg, password_salt, password) VALUES ($1, $2, $3, $4, $5, $6)`,
		id, username, authToken, string(passwordAlgorithm), salt, hashedPassword)
	if err = row.Scan(); isViolatingUniqueConstraintErr(err) {
		return nil, distrybute.ErrUserAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error while inserting new user: %w", err)
	}
	return &distrybute.User{
		ID:                    id,
		Username:              username,
		AuthorizationToken:    authToken,
		PasswordHashAlgorithm: passwordAlgorithm,
	}, nil
}

func (s *service) CheckPassword(username string, password []byte) (ok bool, user *distrybute.User, err error) {
	row := s.connection.QueryRow(context.Background(),
		`SELECT id, username, password, password_alg, password_salt FROM distrybute.users WHERE username LIKE $1`, username)
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
	if isViolatingUniqueConstraintErr(err) {
		return distrybute.ErrUserAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	return nil
}

func (s *service) ResolveAuthorizationToken(id uuid.UUID) (token string, err error) {
	row := s.connection.QueryRow(context.Background(), `SELECT auth_token FROM distrybute.users WHERE id=$1`, id)
	err = row.Scan(&token)
	if err == nil {
		return
	} else if err == pgx.ErrNoRows {
		return "", distrybute.ErrUserNotFound
	} else {
		return "", err
	}
}

func (s *service) RefreshAuthorizationToken(id uuid.UUID) (token string, err error) {
	token, err = generateAuthToken()
	if err != nil {
		return "", err
	}
	row := s.connection.QueryRow(context.Background(), `UPDATE distrybute.users SET auth_token=$1 WHERE id=$2`, token, id)
	err = row.Scan(&token)
	if err == nil {
		return
	} else if isViolatingUniqueConstraintErr(err) {
		return "", distrybute.ErrAuthTokenAlreadyPresent
	} else {
		return "", err
	}
}

func (s *service) DeleteUser(id uuid.UUID) (err error) {
	row := s.connection.QueryRow(context.Background(), `DELETE FROM distrybute.users WHERE id=$1 RETURNING username`, id)
	var username string
	err = row.Scan(&username)
	if err == pgx.ErrNoRows {
		return distrybute.ErrUserNotFound
	} else {
		return
	}
}

func (s *service) UpdatePassword(id uuid.UUID, password []byte) (err error) {
	row := s.connection.QueryRow(context.Background(), `SELECT password_alg, password_salt FROM distrybute.users WHERE id=$1`, id)
	var passwordAlgorithm distrybute.PasswordHashAlgorithm
	var passwordSalt []byte
	err = row.Scan(&passwordAlgorithm, &passwordSalt)
	if err == pgx.ErrNoRows {
		return distrybute.ErrUserNotFound
	} else if err != nil {
		return err
	}
	hashedPassword, err := generatePasswordHash(password, passwordSalt, passwordAlgorithm)
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	row = s.connection.QueryRow(context.Background(),
		`UPDATE distrybute.users SET password=$1 WHERE id=$2`, hashedPassword, id)
	return row.Scan()
}

func isViolatingUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}
	pgErr, ok := err.(*pgconn.PgError)
	// check for unique constraint violation
	return ok && pgErr.Code == "23505"
}
