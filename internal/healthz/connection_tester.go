package healthz

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const postgresCheckTimeout = 2 * time.Second

type ConnectionTester interface {
	TestConnection(ctx context.Context) error
}

type DBTester struct {
	DB *sql.DB
}

func (s *DBTester) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, postgresCheckTimeout)
	defer cancel()
	if err := s.DB.PingContext(ctx); err != nil {
		return fmt.Errorf("test db connection failed: %w", err)
	}
	return nil
}

type AmqpTester struct {
	Conn *amqp.Connection
}

func (s *AmqpTester) TestConnection(ctx context.Context) error {
	if s.Conn.IsClosed() {
		return fmt.Errorf("test amqp connection failed")
	}
	return nil
}

type ChainedTester struct {
	Testers []ConnectionTester
}

func (s *ChainedTester) TestConnection(ctx context.Context) error {
	for _, tester := range s.Testers {
		err := tester.TestConnection(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
