package healthz

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// The Check timeout was set to 2 sec as it seems more than enough time to the DB to answer
const DBCheckTimeout = 2 * time.Second

// ConnectionTester has the method TestConnection that may be used for any external service dependency
type ConnectionTester interface {
	TestConnection(ctx context.Context) error
}

// AmqpTester is the Database ConnectionTester implementation
type DBTester struct {
	DB *sql.DB
}

func (s *DBTester) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, DBCheckTimeout)
	defer cancel()
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("test db connection failed: %w", err)
	}
	return nil
}

// AmqpTester is the RabbitMQ ConnectionTester implementation
type AmqpTester struct {
	Conn *amqp.Connection
}

func (s *AmqpTester) TestConnection(ctx context.Context) error {
	if s.Conn.IsClosed() {
		return fmt.Errorf("test amqp connection failed")
	}
	return nil
}

// ChainedTester stores all ConnectionTesters in a list
type ChainedTester struct {
	Testers []ConnectionTester
}

// TestConnection is a method from ChainedTester that executes all ConnectionTesters in its list
func (s *ChainedTester) TestConnection(ctx context.Context) error {
	for _, tester := range s.Testers {
		err := tester.TestConnection(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

// Implement testing methods dedicated in configuration for each dependency.
