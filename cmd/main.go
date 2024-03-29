package main

import (
	"context"
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
	"github.com/somatom98/badges/user"
	"go.mongodb.org/mongo-driver/mongo"
)

var httpsSrv *http.Server
var conf config.Config
var router *chi.Mux
var mongoDB *mongo.Database

var eventRepository domain.EventRepository
var userRepository domain.UserRepository

var eventConsumer domain.EventConsumer

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

	mongoDB, err = conf.GetMongoDb(context.Background())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to connect to MongoDB")
	}
	log.Info().Msg("MongoDB connected")

	eventRepository = event.NewMongoEventRepository(mongoDB)
	userRepository = user.NewMongoUserRepository(mongoDB)

	eventConsumer = event.NewEventKafkaConsumer(conf.KafkaOptions)

	eventService = event.NewEventService(eventRepository, userRepository, eventConsumer)

	err = eventService.ListenToUserEvents(context.Background())
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Event consumer failed")
	}

	resolver := &graph.Resolver{
		EventService:  eventService,
		EventConsumer: eventConsumer,
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
