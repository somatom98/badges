package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/somatom98/badges/domain"
	"github.com/somatom98/badges/graph/model"
)

// Type is the resolver for the type field.
func (r *eventResolver) Type(ctx context.Context, obj *domain.Event) (string, error) {
	return string(obj.Type), nil
}

// Date is the resolver for the date field.
func (r *eventResolver) Date(ctx context.Context, obj *domain.Event) (string, error) {
	return obj.Date.String(), nil
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, id string) (*domain.EventsList, error) {
	events, err := r.EventService.GetEventsByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

// ReportsEvents is the resolver for the reportsEvents field.
func (r *queryResolver) ReportsEvents(ctx context.Context, mid string) (*domain.EventsList, error) {
	events, err := r.EventService.GetEventsByManagerID(ctx, mid)
	if err != nil {
		return nil, err
	}

	return &events, nil
}

// Events is the resolver for the events field.
func (r *subscriptionResolver) Events(ctx context.Context, id string) (<-chan *domain.Event, error) {
	out := make(chan *domain.Event)

	in, err := r.EventConsumer.Consume(ctx, nil, nil)
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(out)

		for event := range in {
			if event == nil {
				break
			}

			if event.UID == id {
				out <- event
			}
		}
	}()

	return out, nil
}

// Event returns EventResolver implementation.
func (r *Resolver) Event() EventResolver { return &eventResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

type eventResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*domain.Event, error) {
	event := domain.Event{
		ID:   uuid.New().String(),
		UID:  input.User,
		Type: input.Type,
		Date: time.Now().UTC(),
	}

	err := r.EventService.AddUserEvent(ctx, event)
	return &event, err
}
func (r *newEventResolver) Type(ctx context.Context, obj *model.NewEvent, data string) error {
	switch data {
	case string(domain.EventTypeIn):
		obj.Type = domain.EventTypeIn
	case string(domain.EventTypeOut):
		obj.Type = domain.EventTypeOut
	default:
		return fmt.Errorf("err_invalid_value")
	}

	return nil
}
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }
func (r *Resolver) NewEvent() NewEventResolver { return &newEventResolver{r} }

type mutationResolver struct{ *Resolver }
type newEventResolver struct{ *Resolver }
