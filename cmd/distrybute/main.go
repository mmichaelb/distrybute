package main

import (
	"github.com/go-chi/chi/v5"
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
	//	apiRouter := rest.NewRouter(log.With().Str("service", "rest").Logger(), nil, nil)
	//	router.Mount("/api/", apiRouter.BuildHttpHandler())
	panic(http.ListenAndServe("127.0.0.1:8080", router))
}
