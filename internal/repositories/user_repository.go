package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/repositories/common"
)

const (
	pageLimit   int = 100 // Is used as the defaul pageLimit
	oneForToken int = 1   // This will be added to the page limit in order to retrieve the NextPageToken
)

// UserRepo implements models.UserRepository
type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

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
	) RETURNING ID, FIRST_NAME, LAST_NAME, NICKNAME, EMAIL, COUNTRY, CREATED_AT, UPDATED_AT`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("createuser prepare context failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		user.FirstName,
		user.LastName,
		user.Nickname,
		user.Password,
		user.Email,
		user.Country,
	)

	createdUser := &models.User{}
	err = row.Scan(
		&createdUser.ID,
		&createdUser.FirstName,
		&createdUser.LastName,
		&createdUser.Nickname,
		&createdUser.Email,
		&createdUser.Country,
		&createdUser.CreatedAt,
		&createdUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("createuser returned no rows: %w", err)
		}
		return nil, fmt.Errorf("createuser failed: %w", err)
	}
	return createdUser, nil
}

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
	RETURNING ID, FIRST_NAME, LAST_NAME, NICKNAME, EMAIL, COUNTRY, CREATED_AT, UPDATED_AT`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("updateuser prepare context failed: %w", err)
	}
	defer stmt.Close()

	row := stmt.QueryRow(
		user.ID,
		user.FirstName,
		user.LastName,
		user.Nickname,
		user.Password,
		user.Email,
		user.Country,
	)

	updatedUser := &models.User{}
	err = row.Scan(
		&updatedUser.ID,
		&updatedUser.FirstName,
		&updatedUser.LastName,
		&updatedUser.Nickname,
		&updatedUser.Email,
		&updatedUser.Country,
		&updatedUser.CreatedAt,
		&updatedUser.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("updateuser returned no rows: %w", err)
		}
		return nil, fmt.Errorf("updateuser failed: %w", err)
	}
	return updatedUser, nil
}

func (s *UserRepo) FindUsers(ctx context.Context, user *models.User, pageToken string, limit int) (*models.UsersResponse, error) {

	if limit < 1 || limit > pageLimit {
		limit = pageLimit
	}

	userID, err := common.DecodeBase64ToUUID(pageToken)
	if err != nil {
		return nil, fmt.Errorf("findUsers pagetoken decoding failed: %w", err)
	}

	template := `SELECT
		ID,
		FIRST_NAME,
		LAST_NAME,
		NICKNAME,
		EMAIL,
		COUNTRY,
		CREATED_AT,
		UPDATED_AT
		FROM U1.USERS
		WHERE ID >= $1
		{{if .FirstName}} AND FIRST_NAME = '{{.FirstName}}' {{end}}
		{{if .LastName}} AND LAST_NAME = '{{.LastName}}' {{end}}
		{{if .Nickname}} AND NICKNAME = '{{.Nickname}}' {{end}}
		{{if .Country}} AND COUNTRY = '{{.Country}}' {{end}}
		{{if .Email}} AND EMAIL = '{{.Email}}' {{end}}
		{{if .ID}} AND ID >= '{{.ID}}' {{end}}
		ORDER BY ID
		LIMIT $2`

	/* The query above is a template and will be filled depending of the
	values that are existing in the user param.
	The alternative https://github.com/Masterminds/squirrel seemed to me too overpower to use in just one query.
	*/
	query, err := common.ProcessTemplate(template, user)
	if err != nil {
		return nil, fmt.Errorf("findUsers templating failed: %w", err)
	}

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("findUsers preparation failed: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(userID, pageLimit+oneForToken)
	if err != nil {
		return nil, fmt.Errorf("findUsers query failed: %w", err)
	}
	defer rows.Close()

	users := []*models.User{}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Nickname,
			&user.Email,
			&user.Country,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, fmt.Errorf("findUsers returned no rows: %w", err)
				}
				return nil, fmt.Errorf("findUsers failed: %w", err)
			}
		}
		users = append(users, &user)
	}

	// The nextPageToken is always a masked last User.ID
	usersClean, pageToken := common.Paginate(users, limit)

	return &models.UsersResponse{
		Users:     usersClean,
		PageToken: pageToken,
	}, nil
}

func (s *UserRepo) RemoveUser(ctx context.Context, id uuid.UUID) (int64, error) {
	query := "DELETE FROM U1.USERS WHERE id = $1"

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		return 0, fmt.Errorf("removeuser prepare context failed: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, fmt.Errorf("removeuser exec context failed: %w", err)
	}
	return result.RowsAffected()
}
