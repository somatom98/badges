package graph

//go:generate go run github.com/99designs/gqlgen generate

import "github.com/somatom98/badges/graph/model"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	users []*model.User
}
