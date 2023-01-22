package httpAPI_test

import (
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/broker"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/logging"
	usecaseMocks "github.com/jerry0420/queue-system/backend/usecase/mocks"
)

func setUpHttpAPITest() (
	customerUseCase *usecaseMocks.CustomerUseCaseInterface,
	sessionUseCase *usecaseMocks.SessionUseCaseInterface,
	storeUseCase *usecaseMocks.StoreUseCaseInterface,
	integrationUseCase *usecaseMocks.IntegrationUseCaseInterface,
	httpAPIDelivery *httpAPI.HttpAPIDelivery,
	router *mux.Router,
	brokerTool *broker.Broker,
) {
	logger := logging.NewLogger([]string{}, true)
	customerUseCase = new(usecaseMocks.CustomerUseCaseInterface)
	sessionUseCase = new(usecaseMocks.SessionUseCaseInterface)
	storeUseCase = new(usecaseMocks.StoreUseCaseInterface)
	integrationUseCase = new(usecaseMocks.IntegrationUseCaseInterface)
	brokerTool = broker.NewBroker(logger)
	httpAPIDelivery = httpAPI.NewHttpAPIDelivery(
		logger,
		customerUseCase,
		sessionUseCase,
		storeUseCase,
		integrationUseCase,
		brokerTool,
		httpAPI.HttpAPIDeliveryConfig{
			StoreDuration:         time.Duration(2 * time.Second),
			TokenDuration:         time.Duration(2 * time.Second),
			PasswordTokenDuration: time.Duration(2 * time.Second),
			Domain:                "http://localhost.com",
		},
	)

	router = mux.NewRouter()
	return customerUseCase, sessionUseCase, storeUseCase, integrationUseCase, httpAPIDelivery, router, brokerTool
}
