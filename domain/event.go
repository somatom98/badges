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
	GetEventsByUserIDs(ctx context.Context, uids ...string) ([]Event, error)
	AddUserEvent(ctx context.Context, event Event) error
}

type EventService interface {
	GetEventsByUserID(ctx context.Context, uid string) ([]Event, error)
	GetEventsByManagerID(ctx context.Context, managerID string) ([]Event, error)
	AddUserEvent(ctx context.Context, event Event) error
	ListenToUserEvents(ctx context.Context) error
}

type EventConsumer interface {
	Consume(ctx context.Context, handler *func(context.Context, Event) error) (<-chan *Event, error)
}
