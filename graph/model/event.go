package model

import "github.com/somatom98/badges/domain"

type NewEvent struct {
	User string           `json:"user"`
	Type domain.EventType `json:"type"`
}
