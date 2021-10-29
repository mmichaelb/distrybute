package main

import (
	"flag"
	"github.com/go-chi/chi/v5"
	_ "github.com/mmichaelb/distrybute/docs"
	"github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"strconv"
)

func main() {
	host := flag.String("host", "", "the host the web server should listen on")
	port := flag.Int("port", 33761, "the port the web server should listen on")
	router := chi.NewRouter()
	router.Handle("/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:33761/swagger/doc.json"),
	))
	addr := net.JoinHostPort(*host, strconv.Itoa(*port))
	if err := http.ListenAndServe(addr, router); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
