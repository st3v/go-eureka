package test

import (
	"log"
	"net/http"
	"os"

	"github.com/emicklei/go-restful"
)

func NewServer(addr string, debug bool) *http.Server {
	if debug {
		restful.TraceLogger(log.New(os.Stdout, "[restful] ", log.LstdFlags|log.Lshortfile))
	}

	container := restful.NewContainer()
	registry := NewRegistry()
	container.Add(registry.service())

	return &http.Server{
		Addr:    addr,
		Handler: container,
	}
}
