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
	"os/signal"
)

var host string
var port int
var realIpHeader string
var logFile, logLevel string
var minioEndpoint, minioId, minioSecret, minioToken, minioBucket, minioObjectPrefix string

const asciiArt = "\n     _  _       _                 _             _        \n    | |(_)     | |               | |           | |       \n  __| | _  ___ | |_  _ __  _   _ | |__   _   _ | |_  ___ \n / _` || |/ __|| __|| '__|| | | || '_ \\ | | | || __|/ _ \\\n| (_| || |\\__ \\| |_ | |   | |_| || |_) || |_| || |_|  __/\n \\__,_||_||___/ \\__||_|    \\__, ||_.__/  \\__,_| \\__|\\___|\n                            __/ |                        \n                           |___/                         \n"

func RunApp() {
	fmt.Print(asciiArt)
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
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		c.String("postgresuser"), c.String("postgrespassword"),
		c.String("postgreshost"), c.Int("postgresport"),
		c.String("postgresdatabase"))
	log.Info().Msg("connecting to postgres database...")
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal().Err(err).Msg("could not connect to postgres database")
	}
	log.Info().Msg("connecting to minio server...")
	minioClient, err := minio.New(c.String("minioEndpoint"), &minio.Options{
		Creds: credentials.NewStaticV4(c.String("minioId"), c.String("minioSecret"), ""),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could connect to minio server")
	}
	log.Info().Msg("initializing postgres/minio service...")
	service := postgresminio.NewService(conn, minioClient, minioBucket, c.String("minioObjectPrefix"))
	if err = service.Init(); err != nil {
		log.Fatal().Err(err).Msg("could not initialize postgres/minio service")
	}
	router := chi.NewRouter()
	if realIpHeader != "" {
		log.Debug().Str("realIpHeader", realIpHeader).Msg("enabling real ip header detection")
		hookRealIpMiddleware(router)
	}
	apiRouter := controller.NewRouter(log.With().Str("service", "rest").Logger(), service, service)
	router.Mount("/api/", apiRouter)
	router.Get(fmt.Sprintf("/v/{%s}", controller.FileRequestShortIdParamName), apiRouter.HandleFileRequest)
	address := fmt.Sprintf("%s:%d", host, port)
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	server := &http.Server{
		Handler: router,
	}
	go func() {
		log.Info().Str("address", address).Msg("starting server process")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("could not listen and serve web app")
			signalChannel <- os.Interrupt
		}
	}()
	<-signalChannel
	log.Info().Msg("received signal to shut down application")
	return nil
}

func hookRealIpMiddleware(router *chi.Mux) {
	router.Middlewares().Handler(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		realIp := request.Header.Get(realIpHeader)
		request.RemoteAddr = realIp
	}))
}
