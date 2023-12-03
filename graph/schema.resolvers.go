package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/somatom98/badges/graph/model"
)

// CreateEvent is the resolver for the createEvent field.
func (r *mutationResolver) CreateEvent(ctx context.Context, input model.NewEvent) (*model.Event, error) {
	event := model.Event{
		ID:   uuid.New().String(),
		Type: input.Type,
		Date: time.Now().UTC().String(),
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
func (r *queryResolver) Events(ctx context.Context, id *string) ([]*model.Event, error) {
	events := []*model.Event{}
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
				Events: []*model.Event{},
			},
			{
				ID:     "2",
				Name:   "Mario",
				Events: []*model.Event{},
			},
			{
				ID:     "3",
				Name:   "Luca",
				Events: []*model.Event{},
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
func (r *subscriptionResolver) Events(ctx context.Context, id *string) (<-chan *model.Event, error) {
	ch := make(chan *model.Event)

	go func() {
		defer close(ch)

		for {
			time.Sleep(1 * time.Second)

			event := &model.Event{
				ID:   uuid.New().String(),
				Type: "IN",
				Date: time.Now().UTC().String(),
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
func (r *userResolver) Events(ctx context.Context, obj *model.User) ([]*model.Event, error) {
	return obj.Events, nil
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

// Subscription returns SubscriptionResolver implementation.
func (r *Resolver) Subscription() SubscriptionResolver { return &subscriptionResolver{r} }

// User returns UserResolver implementation.
func (r *Resolver) User() UserResolver { return &userResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
type userResolver struct{ *Resolver }