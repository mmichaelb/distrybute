package cli

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/mmichaelb/distrybute/internal/util"
	"github.com/mmichaelb/distrybute/pkg/postgresminio"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io"
	"os"
)

func RunApp() {
	prepareLogger()
	app := util.GeneralApp
	app.Name = "distrybute-cli"
	app.Description = "This CLI application can be used to administrate a distrybute application."
	app.Commands = []*cli.Command{
		userCommand,
	}
	app.Flags = util.PostgresFlags
	app.Before = prepareService
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

func prepareService(c *cli.Context) error {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.String("postgresuser"), c.String("postgrespassword"),
		c.String("postgreshost"), c.Int("postgresport"),
		c.String("postgresdatabase"))
	pool, err := pgxpool.Connect(context.Background(), connString)
	if err != nil {
		return errors.Wrap(err, "could not connect to postgres database")
	}
	service = postgresminio.NewService(pool, nil, "distrybute", "file-")
	return nil
}
