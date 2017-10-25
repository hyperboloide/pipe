package service

import (
	"encoding/json"
	"log"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func RouterFromConfig(config json.RawMessage, silent bool) *chi.Mux {
	r := chi.NewRouter()
	if !silent {
		r.Use(middleware.Logger)
	}

	services := []ServiceDefinition{}
	if err := json.Unmarshal(config, &services); err != nil {
		log.Fatal(err)
	} else if len(services) == 0 {
		log.Fatal("configuration should define at least 1 url.")
	}

	for _, srv := range services {
		if srv.Deleter == nil &&
			srv.ReaderPipe == nil &&
			srv.WriterPipe == nil {
			log.Fatalf("service '%s' should define at least a writer, reader or deleter.", srv.URL)
		}
		if !silent {
			log.Printf("registering service with url '%s'", srv.URL)
		}

		SetHandler(r, srv)
	}
	return r
}
