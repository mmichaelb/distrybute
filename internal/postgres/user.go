package postgres

import (
	"bytes"
	"github.com/google/uuid"
	distrybute "github.com/mmichaelb/distrybute/internal"
	"github.com/rs/zerolog/log"
)

func (s *service) initUserDDL() (err error) {
	rows, err := s.connection.Query(`CREATE TABLE distrybute.users (
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
	rows, err := s.connection.Query(`INSERT INTO distrybute.users (id, username, auth_token, password_alg, password_salt, password) 
		VALUES ($1, $2, $3, $4, $5, $6)`, id, username, authToken, string(passwordAlgorithm), salt, hashedPassword)
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
	rows, err := s.connection.Query(`SELECT (id, username, password, password_alg, password_salt) FROM distrybute.users WHERE user LIKE $1`, username)
	if err != nil {
		return false, nil, err
	}
	if !rows.Next() {
		return false, nil, distrybute.ErrUserNotFound
	}
	values, err := rows.Values()
	if err != nil {
		return false, nil, err
	}
	id := values[0].(uuid.UUID)
	username = values[1].(string)
	expectedPasswordHash := values[2].([]byte)
	passwordAlgorithm := distrybute.PasswordHashAlgorithm(values[3].([]byte))
	passwordSalt := values[4].([]byte)
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

func (s *service) UpdateUsername(user *distrybute.User, newUsername string) (err error) {
	panic("implement me")
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
