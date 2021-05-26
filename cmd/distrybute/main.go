package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/mmichaelb/distrybute/internal/web/rest"
	"github.com/rs/zerolog/log"
)

func main() {
	router := chi.NewRouter()
	apiRouter := rest.NewRouter(log.With().Str("rest").Logger(), nil, nil)
	router.Mount("/api/")
}
