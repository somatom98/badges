package main

import (
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"github.com/somatom98/badges/config"
	"github.com/somatom98/badges/domain"
	"github.com/somatom98/badges/event"
	"github.com/somatom98/badges/graph"
)

var httpsSrv *http.Server
var conf config.Config
var router *chi.Mux

var eventRepository domain.EventRepository

var eventService domain.EventService

func main() {
	conf, err := config.GetFromYaml()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to load config")
	}
	log.Info().
		Msg("Config loaded")

	eventRepository = event.NewMockEventRepository()

	eventService = event.NewEventService(eventRepository)

	resolver := &graph.Resolver{
		EventService: eventService,
	}

	router = chi.NewRouter()
	router.Use(middleware.Timeout(60 * time.Second))
	router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))
	router.Handle("/query", srv)

	if conf.Environment == config.EnvironmentProduction {
		httpsSrv = makeHTTPServer()
		httpsSrv.Addr = ":443"

		go func() {
			log.Info().
				Msg("HTTPS Server starting")
			err = httpsSrv.ListenAndServeTLS("../server.crt", "../server.key")
			if err != nil {
				log.Fatal().
					Err(err).
					Msg("Failed to start HTTPS server")
			}
		}()
	}

	log.Info().
		Msg("HTTP Server starting")
	httpSrv := makeHTTPServer()
	httpSrv.Addr = ":8080"
	err = httpSrv.ListenAndServe()
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to start HTTP server")
	}
}

func makeHTTPServer() *http.Server {
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}
}
