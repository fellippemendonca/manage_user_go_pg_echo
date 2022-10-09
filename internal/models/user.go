package models

//go:generate mockgen -destination=../repositories/user_repository_mock.go -package=repositories . UserRepository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// User Main Struct (It is being used for both API and DB in this project. Ideally should be one for each)
type User struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Nickname  string    `json:"nickname"`
	Password  string    `json:"password,omitempty"` // Hidden from API POST/GET/UPDATE response body for safety
	Email     string    `json:"email"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UsersResponse is a paginated response for the method Get all Users
type UsersResponse struct {
	Users     []*User `json:"users"`
	PageToken string  `json:"page_token"`
}

// UserEvent is a paginated response for the method Get all Users
type UserEvent struct {
	Operation string `json:"operation"`
	UserID    string `json:"user_id"`
	User      *User  `json:"user,omitempty"`
}

// UserRepository has all the methods possible to be called for a User entity
type UserRepository interface {
	// CreateUser creates a new User and returns the User with it's new ID
	CreateUser(ctx context.Context, user *User) (*User, error)
	// FindUsers returnds a paginated list of Users, allowing for filtering by certain criteria (e.g. all Users with the country "UK")
	FindUsers(ctx context.Context, user *User, pageToken string, limit int) (*UsersResponse, error)
	// UpdateUser Modifies an existing User and return the user with its new data
	UpdateUser(ctx context.Context, user *User) (*User, error)
	// RemoveUser deletes a user from the database by its ID
	RemoveUser(ctx context.Context, ID uuid.UUID) (int64, error)
}
