package users_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models_mocks "github.com/fellippemendonca/manage_user_go_pg_echo/internal/models/mocks"
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
	mockedRepo := models_mocks.NewMockUserRepository(ctrl)
	s.UserRepository = mockedRepo

	handler := users.Find(s)

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
			name:       "user found",
			inputID:    "904bc695-6b6c-418a-82a0-0acc7a747d46",
			repoCall:   1,
			repoInput:  uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
			repoResult: 1,
			repoErr:    nil,
			httpStatus: http.StatusAccepted,
		},
		{
			name:       "repo error",
			inputID:    "904bc695-6b6c-418a-82a0-0acc7a747d46",
			repoCall:   1,
			repoInput:  uuid.MustParse("904bc695-6b6c-418a-82a0-0acc7a747d46"),
			repoResult: 1,
			repoErr:    errors.New("Generic Error"),
			httpStatus: http.StatusInternalServerError,
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
