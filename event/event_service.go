package event

import (
	"context"

	"github.com/somatom98/badges/domain"
)

type EventService struct {
	eventRepository domain.EventRepository
}

func NewEventService(eventRepository domain.EventRepository) *EventService {
	return &EventService{
		eventRepository: eventRepository,
	}
}

func (s *EventService) GetEventsByUserID(ctx context.Context, uid string) ([]domain.Event, error) {
	return s.eventRepository.GetEventsByUserID(ctx, uid)
}

func (s *EventService) GetEventsByManagerID(ctx context.Context, managerID string) ([]domain.Event, error) {
	uids := []string{}

	// TODO get uids from manager ID

	return s.eventRepository.GetEventsByIDs(ctx, uids...)
}

func (s *EventService) AddUserEvent(ctx context.Context, event domain.Event) error {
	return s.eventRepository.AddUserEvent(ctx, event)
}
