package postgresminio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/mmichaelb/distrybute/pkg"
)

func (s *Service) CreateNewUser(username string, password []byte) (user *pkg.User, err error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	passwordAlgorithm := pkg.LatestPasswordHashAlgorithm
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
		return nil, pkg.ErrUserAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("error while inserting new user: %w", err)
	}
	return &pkg.User{
		ID:                    id,
		Username:              username,
		AuthorizationToken:    authToken,
		PasswordHashAlgorithm: passwordAlgorithm,
	}, nil
}

func (s *Service) CheckPassword(username string, password []byte) (ok bool, user *pkg.User, err error) {
	row := s.connection.QueryRow(context.Background(),
		`SELECT id, username, password, password_alg, password_salt FROM distrybute.users WHERE username ILIKE $1`, username)
	var id uuid.UUID
	var fetchedUsername string
	var expectedPasswordHash, passwordSalt []byte
	var passwordAlgorithm pkg.PasswordHashAlgorithm
	err = row.Scan(&id, &fetchedUsername, &expectedPasswordHash, &passwordAlgorithm, &passwordSalt)
	if err == pgx.ErrNoRows {
		return false, nil, pkg.ErrUserNotFound
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
	return true, &pkg.User{
		ID:                    id,
		Username:              fetchedUsername,
		PasswordHashAlgorithm: passwordAlgorithm,
	}, nil
}

func (s *Service) UpdateUsername(id uuid.UUID, newUsername string) (err error) {
	row := s.connection.QueryRow(context.Background(), `UPDATE distrybute.users SET username=$1 WHERE id=$2`, newUsername, id)
	err = row.Scan()
	if isViolatingUniqueConstraintErr(err) {
		return pkg.ErrUserAlreadyExists
	} else if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	return nil
}

func (s *Service) ResolveAuthorizationToken(id uuid.UUID) (token string, err error) {
	row := s.connection.QueryRow(context.Background(), `SELECT auth_token FROM distrybute.users WHERE id=$1`, id)
	err = row.Scan(&token)
	if err == nil {
		return
	} else if err == pgx.ErrNoRows {
		return "", pkg.ErrUserNotFound
	} else {
		return "", err
	}
}

func (s *Service) RefreshAuthorizationToken(id uuid.UUID) (token string, err error) {
	token, err = generateAuthToken()
	if err != nil {
		return "", err
	}
	row := s.connection.QueryRow(context.Background(), `UPDATE distrybute.users SET auth_token=$1 WHERE id=$2`, token, id)
	err = row.Scan(&token)
	if isViolatingUniqueConstraintErr(err) {
		return "", pkg.ErrAuthTokenAlreadyPresent
	} else if errors.Is(err, pgx.ErrNoRows) {
		return token, nil
	} else {
		return "", err
	}
}

func (s *Service) GetUserByAuthorizationToken(token string) (bool, *pkg.User, error) {
	row := s.connection.QueryRow(context.Background(), `SELECT id, username FROM distrybute.users WHERE auth_token=$1`, token)
	var id uuid.UUID
	var username string
	err := row.Scan(&id, &username)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil, nil
	} else if err != nil {
		return false, nil, err
	}
	return true, &pkg.User{ID: id, Username: username}, nil
}

func (s *Service) DeleteUser(id uuid.UUID) (err error) {
	row := s.connection.QueryRow(context.Background(), `DELETE FROM distrybute.users WHERE id=$1 RETURNING username`, id)
	var username string
	err = row.Scan(&username)
	if err == pgx.ErrNoRows {
		return pkg.ErrUserNotFound
	} else {
		return
	}
}

func (s *Service) UpdatePassword(id uuid.UUID, password []byte) (err error) {
	row := s.connection.QueryRow(context.Background(), `SELECT password_alg, password_salt FROM distrybute.users WHERE id=$1`, id)
	var passwordAlgorithm pkg.PasswordHashAlgorithm
	var passwordSalt []byte
	err = row.Scan(&passwordAlgorithm, &passwordSalt)
	if err == pgx.ErrNoRows {
		return pkg.ErrUserNotFound
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
