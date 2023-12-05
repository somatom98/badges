package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"github.com/somatom98/badges/domain"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserRepository domain.UserRepository
	EventService   domain.EventService
	EventConsumer  domain.EventKafkaConsumer
}
