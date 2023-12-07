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

func NewEventsList(events []Event) EventsList {
	eventsList := EventsList{}

	for _, event := range events {
		year := event.Date.Year()
		month := int(event.Date.Month())
		day := event.Date.Day()

		var existingYear *YearGroup
		for i := range eventsList.Years {
			if eventsList.Years[i].Year == year {
				existingYear = &eventsList.Years[i]
				break
			}
		}

		if existingYear == nil {
			newYear := YearGroup{
				Year:   year,
				Months: []MonthGroup{},
			}
			existingYear = &newYear
			eventsList.Years = append(eventsList.Years, newYear)
		}

		var existingMonth *MonthGroup
		for i := range existingYear.Months {
			if existingYear.Months[i].Month == month {
				existingMonth = &existingYear.Months[i]
				break
			}
		}

		if existingMonth == nil {
			newMonth := MonthGroup{
				Month: month,
				Days:  []DayGroup{},
			}
			existingMonth = &newMonth
			existingYear.Months = append(existingYear.Months, newMonth)
		}

		var existingDay *DayGroup
		for i := range existingMonth.Days {
			if existingMonth.Days[i].Day == day {
				existingDay = &existingMonth.Days[i]
				break
			}
		}

		if existingDay == nil {
			newDay := DayGroup{
				Day:    day,
				Events: []Event{},
			}
			existingDay = &newDay
			existingMonth.Days = append(existingMonth.Days, newDay)
		}

		existingDay.Events = append(existingDay.Events, event)
	}

	return eventsList
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
