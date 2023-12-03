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

// CreateEvent is the resolver for the createEvent field.
func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*domain.Event, error) {
	event := domain.Event{
		ID: uuid.New().String(),
		// Type: input.Type,
		Date: time.Now().UTC(),
	}

	for _, u := range r.users {
		if u.ID == input.User {
			u.Events = append(u.Events, &event)
			return &event, nil
		}
	}

	return &event, fmt.Errorf("err_user_notfound")
}

// Events is the resolver for the events field.
func (r *queryResolver) Events(ctx context.Context, id *string) ([]*domain.Event, error) {
	events := []*domain.Event{}
	for _, u := range r.users {
		if id == nil || u.ID == *id {
			events = append(events, u.Events...)
		}
	}
	return events, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	if r.users == nil {
		r.users = []*model.User{
			{
				ID:     "1",
				Name:   "Tommaso",
				Events: []*domain.Event{},
			},
			{
				ID:     "2",
				Name:   "Mario",
				Events: []*domain.Event{},
			},
			{
				ID:     "3",
				Name:   "Luca",
				Events: []*domain.Event{},
			},
		}
	}

	for _, u := range r.users {
		if u.ID == id {
			return u, nil
		}
	}

	return nil, nil
}

// Events is the resolver for the events field.
func (r *subscriptionResolver) Events(ctx context.Context, id *string) (<-chan *domain.Event, error) {
	ch := make(chan *domain.Event)

	go func() {
		defer close(ch)

		for {
			time.Sleep(1 * time.Second)

			event := &domain.Event{
				ID:   uuid.New().String(),
				Type: domain.EventTypeIn,
				Date: time.Now().UTC(),
			}
			select {
			case <-ctx.Done():
				return
			case ch <- event:
				// Event sent
			}
		}
	}()

	return ch, nil
}

// Events is the resolver for the events field.
func (r *userResolver) Events(ctx context.Context, obj *model.User) ([]*domain.Event, error) {
	return obj.Events, nil
}

// Event returns EventResolver implementation.
func (r *Resolver) Event() EventResolver { return &eventResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type eventResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }
