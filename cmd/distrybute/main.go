package main

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v4"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	postgresminio2 "github.com/mmichaelb/distrybute/pkg/postgresminio"
	"github.com/mmichaelb/distrybute/pkg/rest/controller"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
)

func main() {
	logFile, err := os.Create("log.json")
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(io.MultiWriter(zerolog.ConsoleWriter{
		Out: os.Stdout,
	}, logFile))
	router := chi.NewRouter()
	conn, err := pgx.Connect(context.Background(), "postgres://postgres:postgres@localhost:5432/postgres")
	if err != nil {
		panic(err)
	}
	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds: credentials.NewStaticV4("minioadmin", "minioadmin", ""),
	})
	if err != nil {
		panic(err)
	}
	service := postgresminio2.NewService(conn, minioClient, "distrybute", "file-")
	if err = service.InitDDL(); err != nil {
		panic(err)
	}
	apiRouter := controller.NewRouter(log.With().Str("service", "rest").Logger(), service, service, service)
	router.Mount("/api/", apiRouter)
	router.Get(fmt.Sprintf("/v/{%s}", controller.FileRequestShortIdParamName), apiRouter.HandleFileRequest)
	panic(http.ListenAndServe("127.0.0.1:8080", router))
}
