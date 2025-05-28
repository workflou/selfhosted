package main

import (
	"flag"
	"net/http"
	"selfhosted/router"
)

var (
	addr = flag.String("addr", ":4000", "HTTP server address")
)

func main() {
	r := router.New()

	s := http.Server{
		Addr:    *addr,
		Handler: r,
	}

	s.ListenAndServe()
}
