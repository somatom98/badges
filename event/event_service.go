package event

import (
	"context"

	"github.com/somatom98/badges/domain"
)

type EventService struct {
	eventRepository domain.EventRepository
	userRepository  domain.UserRepository
	eventConsumer   domain.EventConsumer
}

func NewEventService(eventRepository domain.EventRepository, userRepository domain.UserRepository, eventConsumer domain.EventConsumer) *EventService {
	return &EventService{
		eventRepository: eventRepository,
		userRepository:  userRepository,
		eventConsumer:   eventConsumer,
	}
}

func (s *EventService) GetEventsByUserID(ctx context.Context, uid string) (domain.EventsList, error) {
	events, err := s.eventRepository.GetEventsByUserID(ctx, uid)
	if err != nil {
		return nil, err
	}

	return eventsList(events), nil
}

func (s *EventService) GetEventsByManagerID(ctx context.Context, managerID string) (domain.EventsList, error) {
	uids := []string{}

	users, err := s.userRepository.GetUsersByManagerID(ctx, managerID)
	if err != nil {
		return nil, err
	}

	for _, u := range users {
		uids = append(uids, u.ID)
	}

	events, err := s.eventRepository.GetEventsByUserIDs(ctx, uids...)
	if err != nil {
		return nil, err
	}

	return eventsList(events), nil
}

func (s *EventService) AddUserEvent(ctx context.Context, event domain.Event) error {
	return s.eventRepository.AddUserEvent(ctx, event)
}

func (s *EventService) ListenToUserEvents(ctx context.Context) error {
	handler := s.eventRepository.AddUserEvent
	groupID := "event-consumer"
	_, err := s.eventConsumer.Consume(ctx, &groupID, &handler)
	return err
}

func eventsList(events []domain.Event) domain.EventsList {
	eventsList := domain.EventsList{}

	for _, event := range events {
		year := event.Date.Year()
		month := int(event.Date.Month())
		day := event.Date.Day()

		if eventsList[year] == nil {
			eventsList[year] = map[int]map[int][]domain.Event{}
		}

		if eventsList[year][month] == nil {
			eventsList[year][month] = map[int][]domain.Event{}
		}

		if eventsList[year][month][day] == nil {
			eventsList[year][month][day] = []domain.Event{}
		}

		eventsList[year][month][day] = append(eventsList[year][month][day], event)
	}

	return eventsList
}
