package cli

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"github.com/mmichaelb/distrybute/pkg/postgresminio"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

var postgresUser, postgresPassword, postgresHost string
var postgresPort int
var postgresDatabase string

var postgresFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "postgresuser",
		Aliases:     []string{"pu"},
		Value:       "postgres",
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_USER"},
		Destination: &postgresUser,
	},
	&cli.StringFlag{
		Name:        "postgrespassword",
		Aliases:     []string{"pp"},
		Value:       "postgres",
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_PASSWORD"},
		Destination: &postgresPassword,
	},
	&cli.StringFlag{
		Name:        "postgreshost",
		Aliases:     []string{"ph"},
		Value:       "localhost",
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_HOST"},
		Destination: &postgresHost,
	},
	&cli.IntFlag{
		Name:        "postgresport",
		Aliases:     []string{"pu"},
		Value:       5432,
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_PORT"},
		Destination: &postgresPort,
	},
	&cli.StringFlag{
		Name:        "postgresdatabase",
		Aliases:     []string{"pd"},
		Value:       "postgres",
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_DB"},
		Destination: &postgresDatabase,
	},
}

func RunApp() {
	prepareLogger()
	app := &cli.App{
		Name:  "distrybute-cli",
		Usage: "This CLI application can be used to administrate a distrybute application.",
		Commands: []*cli.Command{
			userCommand,
		},
		Flags:  postgresFlags,
		Before: prepareService,
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func prepareLogger() {
	logFile, err := os.Create("cli-log.json")
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}, logFile))
}

func prepareService(_ *cli.Context) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		postgresUser, postgresPassword, postgresHost, postgresPort, postgresDatabase)
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return errors.Wrap(err, "could not connect to postgres database")
	}
	service = postgresminio.NewService(conn, nil, "distrybute", "file-")
	if err = service.InitDDL(); err != nil {
		return errors.Wrap(err, "could not instantiate postgres/minio service")
	}
	return nil
}
