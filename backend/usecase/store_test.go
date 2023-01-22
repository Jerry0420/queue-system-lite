package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB/mocks"
	"github.com/jerry0420/queue-system/backend/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpStoreTest() (
	pgDBStoreRepository *mocks.PgDBStoreRepositoryInterface,
	pgDBTokenRepository *mocks.PgDBTokenRepositoryInterface,
	storeUsecase usecase.StoreUseCaseInterface,
) {
	pgDBStoreRepository = new(mocks.PgDBStoreRepositoryInterface)
	pgDBTokenRepository = new(mocks.PgDBTokenRepositoryInterface)
	logger := logging.NewLogger([]string{}, true)
	storeUsecase = usecase.NewStoreUsecase(
		pgDBStoreRepository,
		pgDBTokenRepository,
		logger,
		usecase.StoreUsecaseConfig{
			Domain: "http://localhost.com",
		},
	)
	return pgDBStoreRepository, pgDBTokenRepository, storeUsecase
}

func TestUpdateStoreDescription(t *testing.T) {
	pgDBStoreRepository, _, storeUsecase := setUpStoreTest()

	mockStore := domain.Store{
		ID:          1,
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "",
		CreatedAt:   time.Now(),
		Timezone:    "Asia/Taipei",
	}
	expectedMockStoreDescription := "description1"
	pgDBStoreRepository.On("UpdateStore", mock.Anything, &mockStore, "description", expectedMockStoreDescription).
		Return(nil).Run(func(args mock.Arguments) {
			store := args.Get(1).(*domain.Store)
			description := args.Get(3).(string)
			store.Description = description
		}).Once()

	err := storeUsecase.UpdateStoreDescription(context.TODO(), expectedMockStoreDescription, &mockStore)
	assert.NoError(t, err)
	assert.Equal(t, expectedMockStoreDescription, mockStore.Description)
}
