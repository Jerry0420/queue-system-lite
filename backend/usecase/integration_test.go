package usecase_test

import (
	"context"
	"testing"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	pgDBMocks "github.com/jerry0420/queue-system/backend/repository/pgDB/mocks"
	"github.com/jerry0420/queue-system/backend/usecase"
	usecaseMocks "github.com/jerry0420/queue-system/backend/usecase/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setUpIntegrationTest() (
	pgDBTx *pgDBMocks.PgDBTxInterface,
	pgDBStoreRepository *pgDBMocks.PgDBStoreRepositoryInterface,
	pgDBSessionRepository *pgDBMocks.PgDBSessionRepositoryInterface,
	pgDBCustomerRepository *pgDBMocks.PgDBCustomerRepositoryInterface,
	pgDBQueueRepository *pgDBMocks.PgDBQueueRepositoryInterface,
	storeUseCase *usecaseMocks.StoreUseCaseInterface,
	integrationUsecase usecase.IntegrationUseCaseInterface,
	pgDB *pgDBMocks.PgDBInterface,
) {
	pgDBTx = new(pgDBMocks.PgDBTxInterface)
	pgDBStoreRepository = new(pgDBMocks.PgDBStoreRepositoryInterface)
	pgDBSessionRepository = new(pgDBMocks.PgDBSessionRepositoryInterface)
	pgDBCustomerRepository = new(pgDBMocks.PgDBCustomerRepositoryInterface)
	pgDBQueueRepository = new(pgDBMocks.PgDBQueueRepositoryInterface)
	storeUseCase = new(usecaseMocks.StoreUseCaseInterface)
	logger := logging.NewLogger([]string{}, true)

	integrationUsecase = usecase.NewIntegrationUsecase(
		pgDBTx,
		pgDBStoreRepository,
		pgDBSessionRepository,
		pgDBCustomerRepository,
		pgDBQueueRepository,
		storeUseCase,
		nil,
		logger,
		usecase.IntegrationUsecaseConfig{
			StoreDuration:         time.Duration(5 * time.Second),
			TokenDuration:         time.Duration(5 * time.Second),
			PasswordTokenDuration: time.Duration(5 * time.Second),
			ContextTimeOut:        time.Duration(5 * time.Second),
			FromEmail:             "test.gmail.com",
		},
	)
	pgDB = new(pgDBMocks.PgDBInterface)
	return pgDBTx, pgDBStoreRepository, pgDBSessionRepository, pgDBCustomerRepository, pgDBQueueRepository, storeUseCase, integrationUsecase, pgDB
}

func TestCreateStore(t *testing.T) {
	pgDBTx, pgDBStoreRepository, _, _, pgDBQueueRepository, storeUseCase, integrationUsecase, pgDB := setUpIntegrationTest()
	mockStore := domain.Store{
		Email:       "email1",
		Password:    "password1",
		Name:        "name1",
		Description: "description1",
		Timezone:    "Asia/Taipei",
	}
	expectedMockStoreID := 1
	mockQueues := []domain.Queue{
		{
			Name: "queue1",
		},
		{
			Name: "queue2",
		},
	}
	expectedMockQueueID1 := 1
	expectedMockQueueID2 := 2

	storeUseCase.On("EncryptPassword", "password1").Return("encryptPassword1", nil).Once()
	pgDBTx.On("BeginTx").Return(pgDB, nil).Once()
	pgDBTx.On("RollbackTx", pgDB).Once()
	pgDBStoreRepository.On("CreateStore", mock.Anything, pgDB, &mockStore).Return(nil).Run(func(args mock.Arguments) {
		store := args.Get(2).(*domain.Store)
		store.ID = expectedMockStoreID
		store.CreatedAt = time.Now()
	}).Once()

	pgDBQueueRepository.On("CreateQueues", mock.Anything, pgDB, expectedMockStoreID, mockQueues).Return(nil).Run(func(args mock.Arguments) {
		queues := args.Get(3).([]domain.Queue)
		queues[0].ID = expectedMockQueueID1
		queues[0].StoreID = expectedMockStoreID
		queues[1].ID = expectedMockQueueID2
		queues[1].StoreID = expectedMockStoreID
	}).Once()
	pgDBTx.On("CommitTx", pgDB).Return(nil).Once()

	err := integrationUsecase.CreateStore(context.TODO(), &mockStore, mockQueues)
	assert.NoError(t, err)
	assert.Equal(t, "encryptPassword1", mockStore.Password)
	assert.Equal(t, expectedMockStoreID, mockStore.ID)
	assert.Equal(t, expectedMockQueueID1, mockQueues[0].ID)
	assert.Equal(t, expectedMockQueueID2, mockQueues[1].ID)
}
