package messages

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
)

// UserRepo implements models.UserRepo
type UserEvents struct {
	queue          amqp.Queue
	channel        *amqp.Channel
	logger         *zap.Logger
	userRepository models.UserRepository
}

func NewUserEvents(logger *zap.Logger, ch *amqp.Channel, userRepo models.UserRepository) *UserEvents {
	q, err := ch.QueueDeclare(
		"users", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		logger.Fatal("failed to declare event queue", zap.Error(err))
	}

	return &UserEvents{
		queue:          q,
		channel:        ch,
		logger:         logger,
		userRepository: userRepo,
	}
}

func (s *UserEvents) sendEvent(ctx context.Context, msg []byte) error {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	err := s.channel.PublishWithContext(ctx,
		"",           // exchange
		s.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg,
		})
	if err != nil {
		s.logger.Fatal("failed to publish event message", zap.Error(err))
		return err
	}

	return nil
}

// Create user method (Interceptable)
func (s *UserEvents) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return result, err
	}

	// Generate Event
	jsonEvent, err := json.Marshal(&models.UserEvent{
		Operation: "create_user",
		UserID:    result.ID.String(),
		User:      result,
	})
	if err != nil {
		s.logger.Error("failed to marshal create_user message", zap.Error(err))
	} else {
		go s.sendEvent(ctx, jsonEvent)
	}

	return result, err
}

// Update user method
func (s *UserEvents) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Bypass to Repo
	result, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return result, err
	}

	// Generate Event
	jsonEvent, err := json.Marshal(&models.UserEvent{
		Operation: "update_user",
		UserID:    user.ID.String(),
		User:      user,
	})
	if err != nil {
		s.logger.Error("failed to marshal update_user message", zap.Error(err))
	} else {
		go s.sendEvent(ctx, jsonEvent)
	}

	return result, err
}

func (s *UserEvents) FindUsers(ctx context.Context, user *models.User, pageToken string, limit int) (*models.UsersResponse, error) {
	return s.userRepository.FindUsers(ctx, user, pageToken, limit)
}

func (s *UserEvents) RemoveUser(ctx context.Context, id uuid.UUID) (int64, error) {
	result, err := s.userRepository.RemoveUser(ctx, id)
	if err != nil {
		return result, err
	}

	// Generate Event
	jsonEvent, err := json.Marshal(&models.UserEvent{
		Operation: "delete_user",
		UserID:    id.String(),
		User:      nil,
	})
	if err != nil {
		s.logger.Error("failed to marshal delete_user message", zap.Error(err))
	} else {
		go s.sendEvent(ctx, jsonEvent)
	}

	return result, err
}
