package users_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/repositories"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/users"
	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestFind(t *testing.T) {
	s := server.NewServer()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s.Logger = zap.NewNop()
	mockedRepo := repositories.NewMockUserRepository(ctrl)
	s.UserRepository = mockedRepo

	handler := users.Find(s)

	e := echo.New()

	tt := []struct {
		name           string
		queryStr       string
		inputUser      *models.User
		inputPageToken string
		repoCall       int
		repoResult     *models.UsersResponse
		repoErr        error
		httpStatus     int
	}{
		{
			name:     "users.Find StatusOK",
			queryStr: "?first_name=J&last_name=T&nickname=JT&email=j.t@email.com&country=US&page_token=ABC&limit=1",
			repoCall: 1,
			inputUser: &models.User{
				FirstName: "J",
				LastName:  "T",
				Nickname:  "JT",
				Email:     "j.t@email.com",
				Country:   "US",
			},
			inputPageToken: "ABC",
			repoResult: &models.UsersResponse{
				Users: []*models.User{{
					ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
					FirstName: "John",
					LastName:  "Tester",
					Nickname:  "JT",
					Password:  "ABC123!",
					Email:     "john.tester@email.com",
					Country:   "US",
					CreatedAt: time.Time{},
					UpdatedAt: time.Time{},
				}},
				PageToken: "ABC",
			},
			repoErr:    nil,
			httpStatus: http.StatusOK,
		},
		{
			name:     "users.Find StatusNotFound",
			queryStr: "?first_name=J&page_token=ABC&limit=1",
			repoCall: 1,
			inputUser: &models.User{
				FirstName: "J",
			},
			inputPageToken: "ABC",
			repoResult: &models.UsersResponse{
				Users:     []*models.User{},
				PageToken: "",
			},
			repoErr:    sql.ErrNoRows,
			httpStatus: http.StatusNotFound,
		},
		{
			name:     "users.Find StatusInternalServerError",
			queryStr: "?first_name=J&page_token=ABC&limit=1",
			repoCall: 1,
			inputUser: &models.User{
				FirstName: "J",
			},
			inputPageToken: "ABC",
			repoResult: &models.UsersResponse{
				Users:     []*models.User{},
				PageToken: "",
			},
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:           "users.Find StatusBadRequest",
			queryStr:       "?page_token=ABC&limit=17x",
			repoCall:       0,
			inputUser:      &models.User{},
			inputPageToken: "ABC",
			repoResult: &models.UsersResponse{
				Users:     []*models.User{},
				PageToken: "",
			},
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api"+test.queryStr, nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/api")

			test.repoResult.Users = append(test.repoResult.Users, test.inputUser)

			// Mocked User Repository
			mockedRepo.EXPECT().FindUsers(c.Request().Context(), test.inputUser, test.inputPageToken, 1).Times(test.repoCall).Return(test.repoResult, test.repoErr)

			// Assertions
			if assert.NoError(t, handler(c)) {
				assert.Equal(t, test.httpStatus, rec.Code)
				if test.repoErr == nil {
					var usersResponse models.UsersResponse
					err := json.Unmarshal(rec.Body.Bytes(), &usersResponse)
					assert.NoError(t, err)
					assert.Equal(t, *test.repoResult, usersResponse)
				}
			}
		})
	}
}
