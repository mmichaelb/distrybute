package util

import (
	"github.com/urfave/cli/v2"
	"time"
)

var PostgresFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "postgresuser",
		Aliases: []string{"pu"},
		Value:   "postgres",
		EnvVars: []string{"DISTRYBUTE_POSTGRES_USER"},
	},
	&cli.StringFlag{
		Name:    "postgrespassword",
		Aliases: []string{"pp"},
		Value:   "postgres",
		EnvVars: []string{"DISTRYBUTE_POSTGRES_PASSWORD"},
	},
	&cli.StringFlag{
		Name:    "postgreshost",
		Aliases: []string{"ph"},
		Value:   "localhost",
		EnvVars: []string{"DISTRYBUTE_POSTGRES_HOST"},
	},
	&cli.UintFlag{
		Name:    "postgresport",
		Aliases: []string{"ppo"},
		Value:   5432,
		EnvVars: []string{"DISTRYBUTE_POSTGRES_PORT"},
	},
	&cli.StringFlag{
		Name:    "postgresdatabase",
		Aliases: []string{"pd"},
		Value:   "postgres",
		EnvVars: []string{"DISTRYBUTE_POSTGRES_DB"},
	},
	&cli.DurationFlag{
		Name:    "postgrestimeout",
		Aliases: []string{"pt"},
		Value:   time.Second * 30,
	},
}
