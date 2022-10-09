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

// NewUserEvents instantiate a UserEvents
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

// UserEvents implements models.UserRepository and works as a wrapper intercepting all CRUD operations from UserRepository
type UserEvents struct {
	queue          amqp.Queue
	channel        *amqp.Channel
	logger         *zap.Logger
	userRepository models.UserRepository
}

// sendEvent is a method responsbible for dispatching events to RabbitMQ
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

// CreateUser is a method from UserEvents sends a create_user every time a User is created successfully in the DB.
func (s *UserEvents) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	result, err := s.userRepository.CreateUser(ctx, user)
	if err != nil {
		return result, err
	}

	// After the user is created successfully the metod generates the event
	jsonEvent, err := json.Marshal(&models.UserEvent{
		Operation: "create_user",
		UserID:    result.ID.String(),
		User:      result,
	})
	if err != nil {
		s.logger.Error("failed to marshal create_user message", zap.Error(err))
	} else {
		go s.sendEvent(ctx, jsonEvent) // We currently don't care if the event is being handled successfully by the RabbitMQ nor the response should wait for its sending.
	}

	return result, err
}

// UpdateUser is a method from UserEvents sends a create_user every time a User is updated successfully in the DB.
func (s *UserEvents) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Bypass to Repo
	result, err := s.userRepository.UpdateUser(ctx, user)
	if err != nil {
		return result, err
	}

	// After the user is updated successfully the metod generates the event
	jsonEvent, err := json.Marshal(&models.UserEvent{
		Operation: "update_user",
		UserID:    user.ID.String(),
		User:      user,
	})
	if err != nil {
		s.logger.Error("failed to marshal update_user message", zap.Error(err))
	} else {
		go s.sendEvent(ctx, jsonEvent) // We currently don't care if the event is being handled successfully by the RabbitMQ nor the response should wait for its sending.
	}

	return result, err
}

// FindUsers is a method from UserEvents that will simply bypass the call to the UserRepository because we are not broadcasting any reading events. (Maybe we do when we have a caching layer)
func (s *UserEvents) FindUsers(ctx context.Context, user *models.User, pageToken string, limit int) (*models.UsersResponse, error) {
	// Bypass directly to UserRepository.FindUsers
	return s.userRepository.FindUsers(ctx, user, pageToken, limit)
}

// RemoveUser is a method from UserEvents sends a create_user every time a User is deleted successfully in the DB.
func (s *UserEvents) RemoveUser(ctx context.Context, id uuid.UUID) (int64, error) {
	result, err := s.userRepository.RemoveUser(ctx, id)
	if err != nil {
		return result, err
	}

	// After the user is deleted successfully the metod generates the event
	jsonEvent, err := json.Marshal(&models.UserEvent{
		Operation: "delete_user",
		UserID:    id.String(),
		User:      nil,
	})
	if err != nil {
		s.logger.Error("failed to marshal delete_user message", zap.Error(err))
	} else {
		go s.sendEvent(ctx, jsonEvent) // We currently don't care if the event is being handled successfully by the RabbitMQ nor the response should wait for its sending.
	}

	return result, err
}
