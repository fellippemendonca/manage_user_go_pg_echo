package models

//go:generate mockgen -destination=../repositories/user_repository_mock.go -package=repositories . UserRepository

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Main User Struct (It is being used for both API and DB in this project. Ideally should be one for each)
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

// UserRepository
type UserRepository interface {
	// Creates a new User
	CreateUser(ctx context.Context, user *User) (*User, error)
	// Return a paginated list of Users, allowing for filtering by certain criteria (e.g. all Users with the country "UK")
	FindUsers(ctx context.Context, user *User, pageToken string, limit int) (*UsersResponse, error)
	// Modify an existing User
	UpdateUser(ctx context.Context, user *User) (*User, error)
	// Remove a User
	RemoveUser(ctx context.Context, ID uuid.UUID) (int64, error)
}
