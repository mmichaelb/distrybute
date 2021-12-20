package postgresminio

import (
	"context"
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

//go:embed migrations/*
var migrations embed.FS

type Service struct {
	pool         *pgxpool.Pool
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

func (s Service) Init() error {
	m, err := s.instantiateMigrateInstance()
	if err != nil {
		return err
	}
	log.Info().Msg("running database migrations...")
	if err := m.Up(); err == migrate.ErrNoChange {
		log.Info().Msg("no database migration to run")
	} else if err != nil {
		return errors.Wrap(err, "could not run database migrations")
	}
	ok, err := s.minioClient.BucketExists(context.Background(), s.bucketName)
	if err != nil {
		return errors.Wrap(err, "could not check if bucket exists")
	} else if !ok {
		err = s.minioClient.MakeBucket(context.Background(), s.bucketName, minio.MakeBucketOptions{})
		if err != nil {
			return errors.Wrap(err, "could not make new bucket")
		}
		log.Info().Str("bucketName", s.bucketName).Msg("created new bucket")
	}
	return nil
}

func (s Service) instantiateMigrateInstance() (*migrate.Migrate, error) {
	db, err := sql.Open("pgx", s.pool.Config().ConnString())
	if err != nil {
		return nil, errors.Wrap(err, "could not connect to postgres database using pgx driver")
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "could not initiate migrate postgres driver")
	}
	source, err := iofs.New(migrations, "migrations")
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

func NewService(pool *pgxpool.Pool, minioClient *minio.Client, bucketName string, objectPrefix string) *Service {
	return &Service{pool: pool, minioClient: minioClient, bucketName: bucketName, objectPrefix: objectPrefix}
}
