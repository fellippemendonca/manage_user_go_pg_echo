package messages

import (
	"context"
	"fmt"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"

	"github.com/google/uuid"
)

// UserRepo implements models.UserRepo
type UserEvents struct {
	userRepository models.UserRepository
}

func NewUserEvents(userRepo models.UserRepository) *UserEvents {
	return &UserEvents{
		userRepository: userRepo,
	}
}

func (s *UserEvents) sendObjectEvent(name string, user *models.User) error {
	if name == "" || user == nil {
		return fmt.Errorf("name or user not defined")
	}
	fmt.Println("sending " + name + " event.")
	return nil
}

func (s *UserEvents) sendMessageEvent(name string) error {
	if name == "" {
		return fmt.Errorf("name or user not defined")
	}
	fmt.Println("sending " + name + " event.")
	return nil
}

func (s *UserEvents) TestConnection(ctx context.Context) error {
	return s.userRepository.TestConnection(ctx)
}

// Create user method
func (s *UserEvents) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return result, err
	}

	// Event
	s.sendObjectEvent("create_user", user)
	return result, err
}

// Update user method
func (s *UserEvents) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return result, err
	}

	// Event
	s.sendObjectEvent("update_user", user)
	return result, err
}

func (s *UserEvents) FindUsers(ctx context.Context, user *models.User, pageToken string, limit int) ([]*models.User, string, error) {
	return s.userRepository.FindUsers(ctx, user, pageToken, limit)
}

func (s *UserEvents) RemoveUser(ctx context.Context, id uuid.UUID) (int64, error) {
	result, err := s.userRepository.RemoveUser(ctx, id)
	if err != nil {
		return result, err
	}

	// Event
	s.sendMessageEvent("delete_user")
	return result, err
}
