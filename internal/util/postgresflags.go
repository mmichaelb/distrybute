package util

import "github.com/urfave/cli/v2"

var PostgresConnectUriFlag = &cli.StringFlag{
	Name:    "postgresconnecturi",
	Aliases: []string{"postgresuri"},
	// postgres://<user>:<password>@<host>:<port>/<database>?<flags>
	Value:   "postgres://postgres:postgres@localhost:5432/postgres",
	EnvVars: []string{"DISTRYBUTE_POSTGRES_CONNECT_URI"},
}
