package app

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mmichaelb/distrybute/internal/util"
	"github.com/mmichaelb/distrybute/pkg/postgresminio"
	"github.com/mmichaelb/distrybute/pkg/rest/controller"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var host string
var port int
var realIpHeader string
var logFile, logLevel string
var minioEndpoint, minioId, minioSecret, minioToken, minioBucket, minioObjectPrefix string
var pool *pgxpool.Pool
var postgresRetries int
var postgresRetriesInterval time.Duration

const asciiArt = "\n     _  _       _                 _             _        \n    | |(_)     | |               | |           | |       \n  __| | _  ___ | |_  _ __  _   _ | |__   _   _ | |_  ___ \n / _` || |/ __|| __|| '__|| | | || '_ \\ | | | || __|/ _ \\\n| (_| || |\\__ \\| |_ | |   | |_| || |_) || |_| || |_|  __/\n \\__,_||_||___/ \\__||_|    \\__, ||_.__/  \\__,_| \\__|\\___|\n                            __/ |                        \n                           |___/                         \n"

func RunApp() {
	fmt.Print(asciiArt)
	app := util.GeneralApp
	app.Name = "distrybute"
	app.Description = "This application can be used to administrate a distrybute application."
	app.Flags = append(appFlags, util.PostgresConnectUriFlag)
	app.Action = start
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}

func start(c *cli.Context) error {
	err := setupLogging()
	if err != nil {
		return err
	}
	log.Info().Str("version", util.Version).Msg("starting distrybute main application")
	config, err := pgxpool.ParseConfig(c.String("postgresconnecturi"))
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse postgres connect uri")
	}
	log.Debug().Int("retries", postgresRetries).Dur("interval", postgresRetriesInterval).
		Msg("starting postgres connect attempts")
	for i := 0; i < postgresRetries; i++ {
		log.Debug().Int("currentAttempt", i).Send()
		log.Info().Msg("connecting to postgres database server...")
		start := time.Now()
		pool, err = pgxpool.ConnectConfig(context.Background(), config)
		if err != nil {
			log.Warn().Err(err).Msg("could not connect to postgres server")
		} else {
			log.Info().Dur("duration", time.Since(start)).Msg("connected to postgresql server")
			break
		}
		if i != postgresRetries-1 {
			log.Debug().Dur("interval", postgresRetriesInterval).Msg("waiting until next connect attempt...")
			time.Sleep(postgresRetriesInterval)
		} else {
			log.Fatal().Int("retries", postgresRetries).Dur("interval", postgresRetriesInterval).
				Msg("failed to connect to the postgres server within the specified range")
		}
	}
	log.Info().Msg("connecting to minio server...")
	minioClient, err := minio.New(c.String("minioEndpoint"), &minio.Options{
		Creds: credentials.NewStaticV4(c.String("minioId"), c.String("minioSecret"), ""),
	})
	if err != nil {
		log.Fatal().Err(err).Msg("could connect to minio server")
	}
	log.Info().Msg("initializing postgres/minio service...")
	service := postgresminio.NewService(pool, minioClient, minioBucket, c.String("minioObjectPrefix"))
	if err = service.Init(); err != nil {
		log.Fatal().Err(err).Msg("could not initialize postgres/minio service")
	}
	log.Debug().Msg("instantiating new chi router")
	router := chi.NewRouter()
	log.Debug().Str("realIpHeader", realIpHeader).Msg("real ip header output")
	if realIpHeader != "" {
		log.Debug().Str("realIpHeader", realIpHeader).Msg("enabling real ip header detection")
		hookRealIpMiddleware(router)
	}
	log.Debug().Msg("instantiating api router")
	apiRouter := controller.NewRouter(log.With().Str("service", "rest").Logger(), service, service)
	router.Mount("/api/", apiRouter)
	router.Get(fmt.Sprintf("/v/{%s}", controller.FileRequestShortIdParamName), apiRouter.HandleFileRequest)
	log.Debug().Msg("creating channel to listen for interrupts")
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt)
	address := fmt.Sprintf("%s:%d", host, port)
	log.Debug().Str("address", address).Msg("built web server address")
	server := &http.Server{
		Handler: router,
		Addr:    address,
	}
	log.Debug().Msg("starting web server process in separate go routine")
	go func() {
		log.Info().Str("address", address).Msg("starting server process")
		err := server.ListenAndServe()
		log.Debug().Err(err).Msg("stopped listening and serving web server")
		if err != nil && err != http.ErrServerClosed {
			log.Err(err).Msg("could not listen and serve web app")
			signalChannel <- os.Interrupt
		}
	}()
	<-signalChannel
	log.Info().Msg("received signal to shut down application")
	return nil
}

func setupLogging() error {
	logFile, err := os.Create(logFile)
	if err != nil {
		return errors.Wrap(err, "could not create logging output file")
	}
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}, logFile))
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return errors.Wrap(err, "could not parse logging level: "+logLevel)
	}
	zerolog.SetGlobalLevel(level)
	return nil
}

func hookRealIpMiddleware(router *chi.Mux) {
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			realIp := request.Header.Get(realIpHeader)
			if realIp == "" {
				log.Warn().Str("remoteAddr", request.RemoteAddr).Str("expectedRealIpHeader", realIp).
					Msg("request contains no valid real ip header")
			} else {
				sourceIp := request.RemoteAddr
				log.Debug().Str("sourceIp", sourceIp).Str("realIpFromHeader", realIp).Msg("manipulating remote addr field")
				request.RemoteAddr = realIp
			}
			log.Debug().Msg("real ip middleware done - handing request over to further handlers")
			handler.ServeHTTP(writer, request)
		})
	})
}
