package postgresminio

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*
var migrations embed.FS

type service struct {
	connection   *pgx.Conn
	minioClient  *minio.Client
	bucketName   string
	objectPrefix string
}

type wrappedLogger struct {
	logger zerolog.Logger
}

func (w wrappedLogger) Printf(format string, v ...interface{}) {
	w.logger.Printf(format, v...)
}

func (w wrappedLogger) Verbose() bool {
	return w.logger.GetLevel() <= zerolog.DebugLevel
}

func (s service) InitDDL() error {
	m, err := s.instantiateMigrateInstance()
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return errors.Wrap(err, "could not run database migrations")
	}
	return nil
}

func (s service) instantiateMigrateInstance() (*migrate.Migrate, error) {
	db, err := sql.Open("pgx", s.connection.Config().ConnString())
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to postgres database using pgx driver")
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "could not initiate migrate postgres driver")
	}
	source, err := iofs.New(migrations, "./*")
	if err != nil {
		return nil, errors.Wrap(err, "could not initiate FS source for database migrations")
	}
	m, err := migrate.NewWithInstance("fs", source, "postgres", driver)
	if err != nil {
		return nil, errors.Wrap(err, "could not create migrate instance using fs source and postgres db")
	}
	m.Log = wrappedLogger{logger: log.Logger}
	return m, nil
}

func NewService(connection *pgx.Conn, minioClient *minio.Client, bucketName string, objectPrefix string) *service {
	return &service{connection: connection, minioClient: minioClient, bucketName: bucketName, objectPrefix: objectPrefix}
}
