package httpAPI

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/broker"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/middleware"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/usecase"
)

type HttpAPIDeliveryConfig struct {
	StoreDuration         time.Duration
	TokenDuration         time.Duration
	PasswordTokenDuration time.Duration
	Domain                string
	IsProdEnv             bool
}

type HttpAPIDelivery struct {
	logger             logging.LoggerTool
	customerUsecase    usecase.CustomerUseCaseInterface
	sessionUsecase     usecase.SessionUseCaseInterface
	storeUsecase       usecase.StoreUseCaseInterface
	integrationUsecase usecase.IntegrationUseCaseInterface
	broker             *broker.Broker
	config             HttpAPIDeliveryConfig
}

func NewHttpAPIDelivery(
	logger logging.LoggerTool,
	customerUsecase usecase.CustomerUseCaseInterface,
	sessionUsecase usecase.SessionUseCaseInterface,
	storeUsecase usecase.StoreUseCaseInterface,
	integrationUsecase usecase.IntegrationUseCaseInterface,
	broker *broker.Broker,
	config HttpAPIDeliveryConfig,
) *HttpAPIDelivery {
	had := &HttpAPIDelivery{logger, customerUsecase, sessionUsecase, storeUsecase, integrationUsecase, broker, config}
	return had
}

func NewHttpAPIRoutes(
	router *mux.Router,
	mw *middleware.Middleware,
	had *HttpAPIDelivery,
) {
	// stores
	router.HandleFunc(
		V_1("/stores"),
		had.OpenStore,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/signin"),
		had.SigninStore,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/token"),
		had.RefreshToken,
	).Methods(http.MethodPut)

	router.Handle(
		V_1("/stores/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.CloseStore)),
	).Methods(http.MethodDelete)

	router.HandleFunc(
		V_1("/stores/password/forgot"),
		had.ForgotPassword,
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}/password"),
		had.UpdatePassword,
	).Methods(http.MethodPatch).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}/sse"),
		had.GetStoreInfoWithSSE,
	).Methods(http.MethodGet) // get method for sse.

	router.HandleFunc(
		V_1("/stores/{id:[0-9]+}"),
		had.GetStoreInfo,
	).Methods(http.MethodGet)

	router.Handle(
		V_1("/stores/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.UpdateStoreDescription)),
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	router.HandleFunc(
		V_1("/routine/stores"),
		had.CloseStorerRoutine,
	).Methods(http.MethodDelete)

	//queues

	// sessions
	router.HandleFunc(
		V_1("/sessions/sse"),
		had.CreateSession,
	).Methods(http.MethodGet) // get method for sse.

	router.Handle(
		V_1("/sessions/{id}"),
		mw.SessionAuthenticationMiddleware(http.HandlerFunc(had.ScannedSession)),
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	//customers
	router.Handle(
		V_1("/customers"),
		mw.SessionAuthenticationMiddleware(http.HandlerFunc(had.CustomersCreate)),
	).Methods(http.MethodPost).Headers("Content-Type", "application/json")

	router.Handle(
		V_1("/customers/{id:[0-9]+}"),
		mw.AuthenticationMiddleware(http.HandlerFunc(had.CustomerUpdate)),
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	// base routes
	// these two routes will just response to the client directly, and will not go into any middleware.
	router.MethodNotAllowedHandler = http.HandlerFunc(had.methodNotAllow)
	router.NotFoundHandler = http.HandlerFunc(had.notFound)

	// for cors preflight
	router.PathPrefix(V_1("")).HandlerFunc(
		had.preflightHandler,
	).Methods(http.MethodOptions)
}
