package users_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/models"
	models_mocks "github.com/fellippemendonca/manage_user_go_pg_echo/internal/models/mocks"
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
	mockedRepo := models_mocks.NewMockUserRepository(ctrl)
	s.UserRepository = mockedRepo

	handler := users.Update(s)

	e := echo.New()

	tt := []struct {
		name       string
		inputUser  string
		repoCall   int
		repoUser   *models.User
		repoErr    error
		httpStatus int
	}{
		{
			name: "user updated",
			inputUser: `{
				"id": "00000000-0000-0000-0000-000000000000",
				"first_name":"John",
				"last_name":"Tester",
				"nickname":"JT",
				"password":"ABC123!",
				"email":"john.tester@email.com",
				"country":"US"
			}`,
			repoCall: 1,
			repoUser: &models.User{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				FirstName: "John",
				LastName:  "Tester",
				Nickname:  "JT",
				Password:  "ABC123!",
				Email:     "john.tester@email.com",
				Country:   "US",
			},
			repoErr:    nil,
			httpStatus: http.StatusOK,
		},
		{
			name: "repo error",
			inputUser: `{
				"id": "00000000-0000-0000-0000-000000000000",
				"first_name":"John",
				"last_name":"Tester",
				"nickname":"JT",
				"password":"ABC123!",
				"email":"john.tester@email.com",
				"country":"US"
			}`,
			repoCall: 1,
			repoUser: &models.User{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				FirstName: "John",
				LastName:  "Tester",
				Nickname:  "JT",
				Password:  "ABC123!",
				Email:     "john.tester@email.com",
				Country:   "US",
			},
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusInternalServerError,
		},
		{
			name: "broken body",
			inputUser: `{
				"id": "00000000-0000-0000-0000-000000000000",
				"first_name
			}`,
			repoCall: 0,
			repoUser: &models.User{
				ID:        uuid.MustParse("00000000-0000-0000-0000-000000000000"),
				FirstName: "John",
				LastName:  "Tester",
				Nickname:  "JT",
				Password:  "ABC123!",
				Email:     "john.tester@email.com",
				Country:   "US",
			},
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
