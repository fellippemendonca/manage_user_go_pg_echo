package users_test

import (
	"bytes"
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

func TestUpdate(t *testing.T) {
	s := server.NewServer()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s.Logger = zap.NewNop()
	mockedRepo := repositories.NewMockUserRepository(ctrl)
	s.UserRepository = mockedRepo

	handler := users.Update(s)

	e := echo.New()

	tt := []struct {
		name       string
		inputID    string
		inputUser  string
		repoCall   int
		repoUser   *models.User
		repoErr    error
		httpStatus int
	}{
		{
			name:    "users.Update StatusOK",
			inputID: "904bc695-6b6c-418a-82a0-0acc7a747d46",
			inputUser: `{
				"id": "904bc695-6b6c-418a-82a0-0acc7a747d46",
				"first_name":"John",
				"last_name":"Tester",
				"nickname":"JT",
				"password":"ABC123!",
				"email":"john.tester@email.com",
				"country":"US"
			}`,
			repoCall: 1,
			repoUser: &models.User{
				ID:        uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
				FirstName: "John",
				LastName:  "Tester",
				Nickname:  "JT",
				Password:  "ABC123!",
				Email:     "john.tester@email.com",
				Country:   "US",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			repoErr:    nil,
			httpStatus: http.StatusOK,
		},
		{
			name:    "users.Update StatusInternalServerError",
			inputID: "904bc695-6b6c-418a-82a0-0acc7a747d46",
			inputUser: `{
				"id": "904bc695-6b6c-418a-82a0-0acc7a747d46",
				"first_name":"John",
				"last_name":"Tester",
				"nickname":"JT",
				"password":"ABC123!",
				"email":"john.tester@email.com",
				"country":"US"
			}`,
			repoCall: 1,
			repoUser: &models.User{
				ID:        uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
				FirstName: "John",
				LastName:  "Tester",
				Nickname:  "JT",
				Password:  "ABC123!",
				Email:     "john.tester@email.com",
				Country:   "US",
				CreatedAt: time.Time{},
				UpdatedAt: time.Time{},
			},
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:    "users.Update StatusBadRequest",
			inputID: "904bc695-6b6c-418a-82a0-0acc7a747d46",
			inputUser: `{
				"id": "904bc695-6b6c-418a-82a0-0acc7a747d46",
				"first_name":"John",
				"last_name":"Tester",
				"nick
			}`,
			repoCall:   0,
			repoUser:   &models.User{},
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/", bytes.NewReader([]byte(test.inputUser)))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.inputID)

			// Mocked User Repository
			mockedRepo.EXPECT().UpdateUser(c.Request().Context(), test.repoUser).Times(test.repoCall).Return(test.repoUser, test.repoErr)

			// Assertions
			if assert.NoError(t, handler(c)) {
				assert.Equal(t, test.httpStatus, rec.Code)
				if test.repoErr == nil {
					var user models.User
					err := json.Unmarshal(rec.Body.Bytes(), &user)
					assert.NoError(t, err)
					assert.Equal(t, *test.repoUser, user)
				}
			}
		})
	}
}
