package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"

	"github.com/google/uuid"
)

const (
	// TODO: 2sec is plenty, but maybe this should be configurable.
	postgresCheckTimeout       = 2 * time.Second
	pageLimit            int32 = 100
	oneForToken          int32 = 1
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
		return fmt.Errorf("database ping failed: %w", err)
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
		log.Fatal(err)
		return nil, err
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
			return nil, fmt.Errorf("CreateUser failed: %w", err)
		}
		return nil, fmt.Errorf("CreateUser failed: %w", err)
	}
	return user, nil
}

func (s *UserRepo) FindUsers(ctx context.Context, user *models.User) ([]*models.User, error) {

	// t, err := template.New("WHERE").Parse("WHERE You have a task named \"{{ .Name}}\" with description: \"{{ .Description}}\"")
	// if err != nil {
	// 	panic(err)
	// }
	// err = t.Execute(os.Stdout, user)
	// if err != nil {
	// 	panic(err)
	// }

	query := `SELECT
		ID,
		FIRST_NAME,
		LAST_NAME,
		NICKNAME,
		PASSWORD,
		EMAIL,
		COUNTRY,
		CREATED_AT,
		UPDATED_AT
		FROM U1.USERS`

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal(err)
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
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (s *UserRepo) RemoveUser(ctx context.Context, id uuid.UUID) (int64, error) {
	query := "DELETE FROM U1.USERS WHERE id = $1"

	stmt, err := s.db.PrepareContext(ctx, query)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
