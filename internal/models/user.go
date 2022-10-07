package models

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User
type User struct {
	ID         uuid.UUID `json:"id"`
	First_name string    `json:"first_name"`
	Last_name  string    `json:"last_name"`
	Nickname   string    `json:"nickname"`
	Password   string    `json:"password"`
	Email      string    `json:"email"`
	Country    string    `json:"country"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time `json:"updated_at"`
}

// 
type UsersResponse struct {
	Users     []*User `json:"users"`
	PageToken string  `json:"page_token"`
}

// UserRepository
type UserRepository interface {
	TestConnection(ctx context.Context) error
	// Creates a new User
	CreateUser(ctx context.Context, user *User) (*User, error)
	// Return a paginated list of Users, allowing for filtering by certain criteria (e.g. all Users with the country "UK")
	FindUsers(ctx context.Context, user *User, pageToken string, limit int) ([]*User, string, error)
	// Modify an existing User
	UpdateUser(ctx context.Context, user *User) (*User, error)
	// Remove a User
	RemoveUser(ctx context.Context, ID uuid.UUID) (int64, error)
}
