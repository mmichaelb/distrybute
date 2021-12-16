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

var host string
var port int
var realIpHeader string
var logFile, logLevel string
var minioEndpoint, minioId, minioSecret, minioToken, minioBucket, minioObjectPrefix string

func RunApp() {
	app := util.GeneralApp
	app.Name = "distrybute"
	app.Description = "This application can be used to administrate a distrybute application."
	app.Flags = append(appFlags, util.PostgresFlags...)
	app.Action = start
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
	if realIpHeader != "" {
		hookRealIpMiddleware(router)
	}
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

func hookRealIpMiddleware(router *chi.Mux) {
	router.Middlewares().Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		realIp := request.Header.Get(realIpHeader)
		request.RemoteAddr = realIp
	}))
}
