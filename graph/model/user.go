package model

import "github.com/somatom98/badges/domain"

type User struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Events []*domain.Event `json:"events"`
}
