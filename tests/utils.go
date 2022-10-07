package tests

import (
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/messages"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
)

type SUT struct {
	messages.UserEvents
	models.UserRepository
}
