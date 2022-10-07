package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/repositories/common"

	"github.com/google/uuid"
)

const (
	// 2sec is plenty, but maybe this should be configurable.
	postgresCheckTimeout     = 2 * time.Second
	pageLimit            int = 100
	oneForToken          int = 1
)

// UserRepo implements models.UserRepo
type UserRepo struct {
	db *sql.DB
}

// var _ port.UserRepo = (*UserRepo)(nil)

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (s *UserRepo) TestConnection(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, postgresCheckTimeout)
	defer cancel()
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("TestConnection failed: %w", err)
	}
	return nil
}

// Create user method
func (s *UserRepo) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `INSERT INTO U1.USERS (
		FIRST_NAME,
		LAST_NAME,
		NICKNAME,
		PASSWORD,
		EMAIL,
		COUNTRY,
		CREATED_AT,
		UPDATED_AT
	) VALUES (
		$1, --FIRST_NAME
		$2, --LAST_NAME
		$3, --NICKNAME
		$4, --PASSWORD
		$5, --EMAIL
		$6, --COUNTRY
		now() AT TIME ZONE 'utc', -- CREATED_AT
		now() AT TIME ZONE 'utc' -- UPDATED_AT
	) RETURNING ID, FIRST_NAME, LAST_NAME, NICKNAME, PASSWORD, EMAIL, COUNTRY, CREATED_AT, UPDATED_AT`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("CreateUser PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		user.First_name,
		user.Last_name,
		user.Nickname,
		user.Password,
		user.Email,
		user.Country,
	)
	err = row.Scan(
		&user.ID,
		&user.First_name,
		&user.Last_name,
		&user.Nickname,
		&user.Password,
		&user.Email,
		&user.Country,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("CreateUser returned no rows: %w", err)
		}
		return nil, fmt.Errorf("CreateUser failed: %w", err)
	}
	return user, nil
}

// Update user method
func (s *UserRepo) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	query := `UPDATE U1.USERS SET
		FIRST_NAME = $2,
		LAST_NAME = $3,
		NICKNAME = $4,
		PASSWORD = $5,
		EMAIL = $6,
		COUNTRY = $7,
		UPDATED_AT = now() AT TIME ZONE 'utc' -- UPDATED_AT
	WHERE ID = $1
	RETURNING ID, FIRST_NAME, LAST_NAME, NICKNAME, PASSWORD, EMAIL, COUNTRY, CREATED_AT, UPDATED_AT`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("UpdateUser PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		user.ID,
		user.First_name,
		user.Last_name,
		user.Nickname,
		user.Password,
		user.Email,
		user.Country,
	)
	err = row.Scan(
		&user.ID,
		&user.First_name,
		&user.Last_name,
		&user.Nickname,
		&user.Password,
		&user.Email,
		&user.Country,
		&user.Created_at,
		&user.Updated_at,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("UpdateUser returned no rows: %w", err)
		}
		return nil, fmt.Errorf("UpdateUser failed: %w", err)
	}
	return user, nil
}

// Find users by the attributes present in the user struct param
func (s *UserRepo) FindUsers(ctx context.Context, user *models.User, pageToken string, limit int) ([]*models.User, string, error) {

	if limit < 1 || limit > pageLimit {
		limit = pageLimit
	}

	userID, err := common.DecodeBase64ToUUID(pageToken)
	if err != nil {
		return nil, "", fmt.Errorf("FindUsers pageToken decoding failed: %w", err)
	}

	template := `SELECT
		ID,
		FIRST_NAME,
		LAST_NAME,
		NICKNAME,
		PASSWORD,
		EMAIL,
		COUNTRY,
		CREATED_AT,
		UPDATED_AT
		FROM U1.USERS
		WHERE ID >= $1
		{{if .First_name}} AND FIRST_NAME = '{{.First_name}}' {{end}}
		{{if .Last_name}} AND LAST_NAME = '{{.Last_name}}' {{end}}
		{{if .Nickname}} AND NICKNAME = '{{.Nickname}}' {{end}}
		{{if .Country}} AND COUNTRY = '{{.Country}}' {{end}}
		{{if .Email}} AND EMAIL = '{{.Email}}' {{end}}
		{{if .ID}} AND ID >= '{{.ID}}' {{end}}
		ORDER BY ID
		LIMIT $2`

	query, err := common.ProcessTemplate(template, user)
	if err != nil {
		return nil, "", fmt.Errorf("FindUsers Templating failed: %w", err)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, "", fmt.Errorf("FindUsers Preparation failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, pageLimit+oneForToken)
	if err != nil {
		return nil, "", fmt.Errorf("FindUsers Query failed: %w", err)
	}
	defer rows.Close()

	var users []*models.User

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.First_name,
			&user.Last_name,
			&user.Nickname,
			&user.Password,
			&user.Email,
			&user.Country,
			&user.Created_at,
			&user.Updated_at,
		); err != nil {
			log.Fatal(err)
			return nil, "", err
		}
		users = append(users, &user)
	}

	usersClean, pageToken := common.Paginate(users, limit)

	return usersClean, pageToken, nil
}

func (s *UserRepo) RemoveUser(ctx context.Context, id uuid.UUID) (int64, error) {
	query := "DELETE FROM U1.USERS WHERE id = $1"

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("RemoveUser PrepareContext failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("RemoveUser ExecContext failed: %w", err)
	}
	return result.RowsAffected()
}
