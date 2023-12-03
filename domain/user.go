package domain

import "context"

type User struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ManagerID *string `json:"manager_id"`
}

type UserRepository interface {
	GetUserByID(ctx context.Context, uid string) (User, error)
	GetUsersByManagerID(ctx context.Context, managerID string) ([]User, error)
}
