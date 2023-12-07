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

type DayGroup struct {
	Day    int     `json:"day"`
	Events []Event `json:"events"`
}

type MonthGroup struct {
	Month int        `json:"month"`
	Days  []DayGroup `json:"days"`
}

type YearGroup struct {
	Year   int          `json:"year"`
	Months []MonthGroup `json:"months"`
}

type EventsList struct {
	Years []YearGroup `json:"years"`
}

type EventRepository interface {
	GetEventsByUserID(ctx context.Context, uid string) ([]Event, error)
	GetEventsByUserIDs(ctx context.Context, uids ...string) ([]Event, error)
	AddUserEvent(ctx context.Context, event Event) error
}

type EventService interface {
	GetEventsByUserID(ctx context.Context, uid string) (EventsList, error)
	GetEventsByManagerID(ctx context.Context, managerID string) (EventsList, error)
	AddUserEvent(ctx context.Context, event Event) error
	ListenToUserEvents(ctx context.Context) error
}

type EventConsumer interface {
	Consume(ctx context.Context, groupID *string, handler *func(context.Context, Event) error) (<-chan *Event, error)
}
