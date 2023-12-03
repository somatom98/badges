package user

import (
	"context"
	"fmt"

	"github.com/somatom98/badges/domain"
)

type MockUserRepository struct {
	users []domain.User
}

func NewMockUserRepository() *MockUserRepository {
	marioManager := "1"
	lucaManager := "2"

	return &MockUserRepository{
		users: []domain.User{
			{
				ID:        "1",
				Name:      "Tommaso",
				ManagerID: nil,
			},
			{
				ID:        "2",
				Name:      "Mario",
				ManagerID: &marioManager,
			},
			{
				ID:        "3",
				Name:      "Luca",
				ManagerID: &lucaManager,
			},
		},
	}
}

func (r *MockUserRepository) GetUserByID(ctx context.Context, uid string) (domain.User, error) {
	for _, u := range r.users {
		if u.ID == uid {
			return u, nil
		}
	}

	return domain.User{}, fmt.Errorf("err_not_found")
}

func (r *MockUserRepository) GetUsersByManagerID(ctx context.Context, managerID string) ([]domain.User, error) {
	users := []domain.User{}

	for _, u := range r.users {
		if u.ManagerID == &managerID {
			users = append(users, u)
		}
	}

	return users, nil
}
