package usecase_test

import (
	"context"
	"testing"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB/mocks"
	"github.com/jerry0420/queue-system/backend/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpSessionTest() (
	pgDBSessionRepository *mocks.PgDBSessionRepositoryInterface,
	sessionUsecase usecase.SessionUseCaseInterface,
) {
	pgDBSessionRepository = new(mocks.PgDBSessionRepositoryInterface)
	logger := logging.NewLogger([]string{}, true)
	sessionUsecase = usecase.NewSessionUsecase(
		pgDBSessionRepository,
		logger,
	)
	return pgDBSessionRepository, sessionUsecase
}

func TestGetSessionById(t *testing.T) {
	pgDBSessionRepository, sessionUsecase := setUpSessionTest()

	t.Run("exist session id", func(t *testing.T) {
		expectedMockSession := domain.StoreSession{
			ID:                 "im_esssion_id",
			StoreId:            1,
			StoreSessionState: domain.StoreSessionState.SCANNED,
		}

		pgDBSessionRepository.
			On("GetSessionById", mock.Anything, expectedMockSession.ID).
			Return(expectedMockSession, nil).
			Once()

		session, err := sessionUsecase.GetSessionById(context.TODO(), expectedMockSession.ID)
		assert.NoError(t, err)
		assert.Equal(t, expectedMockSession.ID, session.ID)
	})

	t.Run("non-exist session id", func(t *testing.T) {
		expectedMockSession := domain.StoreSession{}
		session, err := sessionUsecase.GetSessionById(context.TODO(), expectedMockSession.ID)
		assert.NotNil(t, err)
		assert.Equal(t, expectedMockSession, session)
	})
}
