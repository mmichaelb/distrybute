package app

import (
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"
	"time"
)

var appFlags = []cli.Flag{
	&cli.StringFlag{
		Name:        "host",
		EnvVars:     []string{"DISTRYBUTE_HOST"},
		Value:       "",
		Destination: &host,
	},
	&cli.IntFlag{
		Name:        "port",
		EnvVars:     []string{"DISTRYBUTE_PORT"},
		Value:       10711,
		Destination: &port,
	},
	&cli.StringFlag{
		Name:        "realIpHeader",
		EnvVars:     []string{"DISTRYBUTE_REAL_IP_HEADER"},
		Destination: &realIpHeader,
	},
	&cli.StringFlag{
		Name:        "logFile",
		EnvVars:     []string{"DISTRYBUTE_LOG_FILE"},
		Value:       "log.json",
		Destination: &logFile,
	},
	&cli.StringFlag{
		Name:        "level",
		EnvVars:     []string{"DISTRYBUTE_LEVEL"},
		Value:       zerolog.InfoLevel.String(),
		Destination: &logLevel,
	},
	&cli.StringFlag{
		Name:        "minioEndpoint",
		EnvVars:     []string{"DISTRYBUTE_MINIO_ENDPOINT"},
		Destination: &minioEndpoint,
		Required:    true,
	},
	&cli.StringFlag{
		Name:        "minioId",
		EnvVars:     []string{"DISTRYBUTE_MINIO_ID"},
		Destination: &minioId,
	},
	&cli.StringFlag{
		Name:        "minioSecret",
		EnvVars:     []string{"DISTRYBUTE_MINIO_SECRET"},
		Destination: &minioSecret,
	},
	&cli.StringFlag{
		Name:        "minioToken",
		EnvVars:     []string{"DISTRYBUTE_MINIO_TOKEN"},
		Value:       "",
		Destination: &minioToken,
	},
	&cli.StringFlag{
		Name:        "minioBucket",
		EnvVars:     []string{"DISTRYBUTE_MINIO_BUCKET"},
		Destination: &minioBucket,
		Required:    true,
	},
	&cli.StringFlag{
		Name:        "minioObjectPrefix",
		EnvVars:     []string{"DISTRYBUTE_MINIO_OBJECT_PREFIX"},
		Value:       "file-",
		Destination: &minioObjectPrefix,
	},
	&cli.IntFlag{
		Name:        "postgresretries",
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_RETRIES"},
		Value:       12,
		Destination: &postgresRetries,
	},
	&cli.DurationFlag{
		Name:        "postgresretryinterval",
		EnvVars:     []string{"DISTRYBUTE_POSTGRES_RETRIES_INTERVAL"},
		Value:       time.Second * 5,
		Destination: &postgresRetriesInterval,
	},
}
