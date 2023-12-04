package event

import (
	"context"

	"github.com/somatom98/badges/domain"
)

type EventService struct {
	eventRepository domain.EventRepository
	userRepository  domain.UserRepository
}

func NewEventService(eventRepository domain.EventRepository, userRepository domain.UserRepository) *EventService {
	return &EventService{
		eventRepository: eventRepository,
		userRepository:  userRepository,
	}
}

func (s *EventService) GetEventsByUserID(ctx context.Context, uid string) ([]domain.Event, error) {
	return s.eventRepository.GetEventsByUserID(ctx, uid)
}

func (s *EventService) GetEventsByManagerID(ctx context.Context, managerID string) ([]domain.Event, error) {
	uids := []string{}

	users, err := s.userRepository.GetUsersByManagerID(ctx, managerID)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		uids = append(uids, u.ID)
	}

	return s.eventRepository.GetEventsByUserIDs(ctx, uids...)
}

func (s *EventService) AddUserEvent(ctx context.Context, event domain.Event) error {
	return s.eventRepository.AddUserEvent(ctx, event)
}
