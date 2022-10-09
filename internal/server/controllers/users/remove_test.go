package users_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/repositories"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server"
	"github.com/fellippemendonca/manage_user_go_pg_echo/internal/server/controllers/users"
	"github.com/google/uuid"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestRemove(t *testing.T) {
	s := server.NewServer()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s.Logger = zap.NewNop()
	mockedRepo := repositories.NewMockUserRepository(ctrl)
	s.UserRepository = mockedRepo

	handler := users.Remove(s)

	e := echo.New()

	tt := []struct {
		name       string
		inputID    string
		repoInput  uuid.UUID
		repoCall   int
		repoResult int64
		repoErr    error
		httpStatus int
	}{
		{
			name:       "users.Remove StatusAccepted",
			inputID:    "904bc695-6b6c-418a-82a0-0acc7a747d46",
			repoCall:   1,
			repoInput:  uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
			repoResult: 1,
			repoErr:    nil,
			httpStatus: http.StatusAccepted,
		},
		{
			name:       "users.Remove StatusInternalServerError",
			inputID:    "904bc695-6b6c-418a-82a0-0acc7a747d46",
			repoCall:   1,
			repoInput:  uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
			repoResult: 1,
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusInternalServerError,
		},
		{
			name:       "users.Remove StatusBadRequest",
			inputID:    "904bc695-6b6cxxxxx82a0-0acc7a747d46",
			repoCall:   0,
			repoInput:  uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
			repoResult: 0,
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusBadRequest,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodDelete, "/", nil)
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/:id")
			c.SetParamNames("id")
			c.SetParamValues(test.inputID)

			// Mocked User Repository
			mockedRepo.EXPECT().RemoveUser(c.Request().Context(), test.repoInput).Times(test.repoCall).Return(test.repoResult, test.repoErr)

			// Assertions
			if assert.NoError(t, handler(c)) {
				assert.Equal(t, test.httpStatus, rec.Code)
			}
		})
	}
}
