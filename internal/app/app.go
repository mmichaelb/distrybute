package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mmichaelb/distrybute/internal/util"
	"github.com/mmichaelb/distrybute/pkg/postgresminio"
	"github.com/mmichaelb/distrybute/pkg/rest/controller"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
)

var (
	GitBranch    string
	GitTag       string
	GitCommitSha string
)

var host string
var port int
var logFile, logLevel string
var minioEndpoint, minioId, minioSecret, minioToken, minioBucket, minioObjectPrefix string

func RunApp() {
	app := &cli.App{
		Name:    "distrybute",
		Usage:   "This application can be used to administrate a distrybute application.",
		Authors: []*cli.Author{{Name: "mmichaelb", Email: "me@mmichaelb.pw"}},
		Version: fmt.Sprintf("%s/%s/%s", GitBranch, GitTag, GitCommitSha),
		Flags:   append(appFlags, util.PostgresFlags...),
		Action:  start,
	}
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func start(c *cli.Context) error {
	logFile, err := os.Create(logFile)
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}, logFile))
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		panic(err)
	}
	log.Level(level)
	router := chi.NewRouter()
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.String("postgresuser"), c.String("postgrespassword"),
		c.String("postgreshost"), c.Int("postgresport"),
		c.String("postgresdatabase"))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		panic(err)
	}
	minioClient, err := minio.New(c.String("minioEndpoint"), &minio.Options{
		Creds: credentials.NewStaticV4(c.String("minioId"), c.String("minioSecret"), ""),
	})
	if err != nil {
		panic(err)
	}
	service := postgresminio.NewService(conn, minioClient, minioBucket, c.String("minioObjectPrefix"))
	if err = service.Init(); err != nil {
		panic(err)
	}
	apiRouter := controller.NewRouter(log.With().Str("service", "rest").Logger(), service, service)
	router.Mount("/api/", apiRouter)
	router.Get(fmt.Sprintf("/v/{%s}", controller.FileRequestShortIdParamName), apiRouter.HandleFileRequest)
	panic(http.ListenAndServe(fmt.Sprintf("%s:%d", c.String("host"), c.Int("port")), router))
}
