package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/jerry0420/queue-system/backend/broker"
	"github.com/jerry0420/queue-system/backend/config"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/middleware"
	"github.com/jerry0420/queue-system/backend/logging"

	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/jerry0420/queue-system/backend/usecase"
	"github.com/jerry0420/queue-system/backend/utils"
)

func main() {
	logger := logging.NewLogger([]string{"method", "url", "code", "sep", "requestID", "duration"}, false)

	err := utils.InitDirPath("csvs")
	if err != nil {
		log.Fatalf("%v", err)
	}

	var db *sql.DB
	dbLocation := config.ServerConfig.POSTGRES_LOCATION()

	db = pgDB.GetDb(config.ServerConfig.POSTGRES_USER(), config.ServerConfig.POSTGRES_PASSWORD(), dbLocation, logger)

	defer func() {
		err := db.Close()
		if err != nil {
			logger.ERRORf("db connection close fail %v", err)
		}
	}()

	router := mux.NewRouter()

	pgDBTx := pgDB.NewPgDBTx(db, logger)
	pgDBStoreRepository := pgDB.NewPgDBStoreRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBTokenRepository := pgDB.NewPgDBTokenRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBSessionRepository := pgDB.NewPgDBSessionRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBQueueRepository := pgDB.NewPgDBQueueRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())
	pgDBCustomerRepository := pgDB.NewPgDBCustomerRepository(db, logger, config.ServerConfig.CONTEXT_TIMEOUT())

	storeUsecase := usecase.NewStoreUsecase(
		pgDBStoreRepository,
		pgDBTokenRepository,
		logger,
		usecase.StoreUsecaseConfig{
			Domain:       config.ServerConfig.DOMAIN(),
			TokenSignKey: config.ServerConfig.TOKEN_SIGN_KEY(),
		},
	)
	sessionUsecase := usecase.NewSessionUsecase(pgDBSessionRepository, logger)
	customerUsecase := usecase.NewCustomerUsecase(pgDBCustomerRepository, logger)
	dialer := utils.InitEmailDialer(
		config.ServerConfig.EMAIL_SERVER(),
		config.ServerConfig.EMAIL_PORT(),
		config.ServerConfig.EMAIL_USERNAME(),
		config.ServerConfig.EMAIL_PASSWORD(),
	)
	integrationUsecase := usecase.NewIntegrationUsecase(
		pgDBTx,
		pgDBStoreRepository,
		pgDBSessionRepository,
		pgDBCustomerRepository,
		pgDBQueueRepository,
		storeUsecase,
		dialer,
		logger,
		usecase.IntegrationUsecaseConfig{
			StoreDuration:         config.ServerConfig.STOREDURATION(),
			TokenDuration:         config.ServerConfig.TOKENDURATION(),
			PasswordTokenDuration: config.ServerConfig.PASSWORDTOKENDURATION(),
			ContextTimeOut:        config.ServerConfig.CONTEXT_TIMEOUT() * 4,
			FromEmail:             config.ServerConfig.EMAIL_FROM(),
		},
	)

	broker := broker.NewBroker(logger)
	defer broker.CloseAll()

	httpAPIDelivery := httpAPI.NewHttpAPIDelivery(
		logger,
		customerUsecase,
		sessionUsecase,
		storeUsecase,
		integrationUsecase,
		broker,
		httpAPI.HttpAPIDeliveryConfig{
			StoreDuration:         config.ServerConfig.STOREDURATION(),
			TokenDuration:         config.ServerConfig.TOKENDURATION(),
			PasswordTokenDuration: config.ServerConfig.PASSWORDTOKENDURATION(),
			Domain:                config.ServerConfig.DOMAIN(),
			IsProdEnv:             config.ServerConfig.ENV() == config.EnvState.PROD,
		},
	)

	mw := middleware.NewMiddleware(router, logger, integrationUsecase, sessionUsecase)

	httpAPI.NewHttpAPIRoutes(router, mw, httpAPIDelivery)

	server := &http.Server{
		Addr:         "0.0.0.0:8000",
		WriteTimeout: time.Hour * 24,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router,
	}

	go func() {
		logger.INFOf("Server Start!")
		if err := server.ListenAndServe(); err != nil {
			logger.ERRORf("ListenAndServe http fail %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until receive signal...
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	server.Shutdown(ctx)
	logger.INFOf("shutting down")
	os.Exit(0)
}
