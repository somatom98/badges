package domain

import (
	"context"
	"time"
)

type Event struct {
	ID   string    `json:"id"`
	UID  string    `json:"uid"`
	Type EventType `json:"type"`
	Date time.Time `json:"date"`
}

type EventType string

const (
	EventTypeIn  EventType = "IN"
	EventTypeOut EventType = "OUT"
)

type EventRepository interface {
	GetEventsByUserID(ctx context.Context, uid string) ([]Event, error)
	GetEventsByIDs(ctx context.Context, uids ...string) ([]Event, error)
	AddUserEvent(ctx context.Context, event Event) error
}
