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
		return domain.EventsList{}, err
	}

	return eventsList(events), nil
}

func (s *EventService) GetEventsByManagerID(ctx context.Context, managerID string) (domain.EventsList, error) {
	uids := []string{}

	users, err := s.userRepository.GetUsersByManagerID(ctx, managerID)
	if err != nil {
		return domain.EventsList{}, err
	}

	for _, u := range users {
		uids = append(uids, u.ID)
	}

	events, err := s.eventRepository.GetEventsByUserIDs(ctx, uids...)
	if err != nil {
		return domain.EventsList{}, err
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

		var existingYear *domain.YearGroup
		for i := range eventsList.Years {
			if eventsList.Years[i].Year == year {
				existingYear = &eventsList.Years[i]
				break
			}
		}

		if existingYear == nil {
			newYear := domain.YearGroup{
				Year:   year,
				Months: []domain.MonthGroup{},
			}
			existingYear = &newYear
			eventsList.Years = append(eventsList.Years, newYear)
		}

		var existingMonth *domain.MonthGroup
		for i := range existingYear.Months {
			if existingYear.Months[i].Month == month {
				existingMonth = &existingYear.Months[i]
				break
			}
		}

		if existingMonth == nil {
			newMonth := domain.MonthGroup{
				Month: month,
				Days:  []domain.DayGroup{},
			}
			existingMonth = &newMonth
			existingYear.Months = append(existingYear.Months, newMonth)
		}

		var existingDay *domain.DayGroup
		for i := range existingMonth.Days {
			if existingMonth.Days[i].Day == day {
				existingDay = &existingMonth.Days[i]
				break
			}
		}

		if existingDay == nil {
			newDay := domain.DayGroup{
				Day:    day,
				Events: []domain.Event{},
			}
			existingDay = &newDay
			existingMonth.Days = append(existingMonth.Days, newDay)
		}

		existingDay.Events = append(existingDay.Events, event)
	}

	return eventsList
}
