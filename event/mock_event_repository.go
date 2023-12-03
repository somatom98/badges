package event

import (
	"context"

	"github.com/somatom98/badges/domain"
)

type MockEventRepository struct {
	events []domain.Event
}

func NewMockEventRepository() *MockEventRepository {
	return &MockEventRepository{
		events: []domain.Event{},
	}
}

func (r *MockEventRepository) GetEventsByUserID(ctx context.Context, uid string) ([]domain.Event, error) {
	filteredEvents := []domain.Event{}

	for _, e := range r.events {
		if e.UID == uid {
			filteredEvents = append(filteredEvents, e)
		}
	}

	return filteredEvents, nil
}

func (r *MockEventRepository) GetEventsByIDs(ctx context.Context, uids ...string) ([]domain.Event, error) {
	filteredEvents := []domain.Event{}

	for _, e := range r.events {
		for _, uid := range uids {
			if e.UID == uid {
				filteredEvents = append(filteredEvents, e)
			}
			continue
		}
	}

	return filteredEvents, nil
}

func (r *MockEventRepository) AddUserEvent(ctx context.Context, event domain.Event) error {
	r.events = append(r.events, event)
	return nil
}
