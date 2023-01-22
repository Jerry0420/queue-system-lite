package usecase

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	gomail "gopkg.in/gomail.v2"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/jerry0420/queue-system/backend/utils"
)

type IntegrationUsecaseConfig struct {
	StoreDuration         time.Duration
	TokenDuration         time.Duration
	PasswordTokenDuration time.Duration
	ContextTimeOut        time.Duration
	FromEmail             string
}

type integrationUsecase struct {
	pgDBTx                 pgDB.PgDBTxInterface
	pgDBStoreRepository    pgDB.PgDBStoreRepositoryInterface
	pgDBSessionRepository  pgDB.PgDBSessionRepositoryInterface
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface
	pgDBQueueRepository    pgDB.PgDBQueueRepositoryInterface
	storeUsecase           StoreUseCaseInterface
	dialer                 *gomail.Dialer
	logger                 logging.LoggerTool
	config                 IntegrationUsecaseConfig
}

func NewIntegrationUsecase(
	pgDBTx pgDB.PgDBTxInterface,
	pgDBStoreRepository pgDB.PgDBStoreRepositoryInterface,
	pgDBSessionRepository pgDB.PgDBSessionRepositoryInterface,
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface,
	pgDBQueueRepository pgDB.PgDBQueueRepositoryInterface,
	storeUsecase StoreUseCaseInterface,
	dialer *gomail.Dialer,
	logger logging.LoggerTool,
	config IntegrationUsecaseConfig,
) IntegrationUseCaseInterface {
	return &integrationUsecase{
		pgDBTx,
		pgDBStoreRepository,
		pgDBSessionRepository,
		pgDBCustomerRepository,
		pgDBQueueRepository,
		storeUsecase,
		dialer,
		logger,
		config,
	}
}

func (iu *integrationUsecase) CreateCustomers(
	ctx context.Context,
	session *domain.StoreSession,
	oldState string,
	newState string,
	customers []domain.Customer,
) error {
	tx, err := iu.pgDBTx.BeginTx()
	if err != nil {
		return err
	}
	defer iu.pgDBTx.RollbackTx(tx)

	err = iu.pgDBSessionRepository.UpdateSessionState(ctx, tx, session, oldState, newState)
	if err != nil {
		return err
	}
	session.StoreSessionState = newState

	err = iu.pgDBCustomerRepository.CreateCustomers(ctx, tx, customers)
	if err != nil {
		return err
	}

	err = iu.pgDBTx.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func (iu *integrationUsecase) CreateStore(ctx context.Context, store *domain.Store, queues []domain.Queue) error {
	encryptedPassword, err := iu.storeUsecase.EncryptPassword(store.Password)
	if err != nil {
		return err
	}
	store.Password = encryptedPassword

	tx, err := iu.pgDBTx.BeginTx()
	if err != nil {
		return err
	}
	defer iu.pgDBTx.RollbackTx(tx)

	err = iu.pgDBStoreRepository.CreateStore(ctx, tx, store)
	if err != nil {
		return err
	}

	err = iu.pgDBQueueRepository.CreateQueues(ctx, tx, store.ID, queues)
	if err != nil {
		return err
	}

	err = iu.pgDBTx.CommitTx(tx)
	if err != nil {
		return err
	}
	return nil
}

func (iu *integrationUsecase) SigninStore(ctx context.Context, email string, password string) (store domain.Store, token string, refreshTokenExpiresAt time.Time, err error) {
	store, err = iu.pgDBStoreRepository.GetStoreByEmail(ctx, email)
	err = iu.storeUsecase.ValidatePassword(store.Password, password)
	if err != nil {
		return store, token, refreshTokenExpiresAt, err
	}

	// let crontab take responsibility of "closestore" tasks.
	refreshTokenExpiresAt = time.Now().Add(iu.config.StoreDuration)
	// refreshTokenExpiresAt = store.CreatedAt.Add(uc.config.StoreDuration)
	token, err = iu.storeUsecase.GenerateToken(
		ctx,
		store,
		domain.TokenTypes.REFRESH,
		refreshTokenExpiresAt,
	)
	if err != nil {
		return store, token, refreshTokenExpiresAt, err
	}

	return store, token, refreshTokenExpiresAt, nil
}

func (iu *integrationUsecase) RefreshToken(ctx context.Context, encryptedRefreshToken string) (
	store domain.Store,
	normalToken string,
	sessionToken string,
	tokenExpiresAt time.Time,
	err error,
) {
	tokenClaims, err := iu.storeUsecase.VerifyToken(
		ctx,
		encryptedRefreshToken,
		domain.TokenTypes.REFRESH,
		true,
	)
	if err != nil {
		return store, normalToken, sessionToken, tokenExpiresAt, err
	}
	store = domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}

	tokenExpiresAt = time.Now().Add(iu.config.TokenDuration)
	// normal token
	normalToken, err = iu.storeUsecase.GenerateToken(
		ctx,
		store,
		domain.TokenTypes.NORMAL,
		tokenExpiresAt,
	)
	if err != nil {
		return store, normalToken, sessionToken, tokenExpiresAt, err
	}
	// session token
	sessionToken, err = iu.storeUsecase.GenerateToken(
		ctx,
		store,
		domain.TokenTypes.SESSION,
		tokenExpiresAt,
	)
	if err != nil {
		return store, normalToken, sessionToken, tokenExpiresAt, err
	}

	return store, normalToken, sessionToken, tokenExpiresAt, nil
}

func (iu *integrationUsecase) CloseStore(ctx context.Context, store domain.Store) error {
	tx, err := iu.pgDBTx.BeginTx()
	if err != nil {
		return err
	}
	defer iu.pgDBTx.RollbackTx(tx)

	customers, err := iu.pgDBCustomerRepository.GetCustomersWithQueuesByStore(ctx, tx, &store)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(1)
	go func() {
		defer wg.Done()
		date, csvFileName, csvContent := iu.storeUsecase.GenerateCsvFileNameAndContent(store.CreatedAt, store.Timezone, store.Name, customers)
		filePath, err := utils.GenerateCSV(
			ctx,
			iu.config.ContextTimeOut*4,
			csvFileName,
			csvContent,
		)
		if err != nil {
			return
		}
		emailSubject, emailContent := iu.storeUsecase.GenerateEmailContentOfCloseStore(store.Name, date)
		_, err = utils.SendEmail(
			ctx,
			iu.config.ContextTimeOut*4,
			iu.dialer,
			iu.config.FromEmail,
			emailSubject,
			emailContent,
			store.Email,
			filePath,
		)
	}()

	wg.Add(1)
	go func(errChan chan error) {
		defer wg.Done()
		// TODO: open it later...
		err := iu.pgDBStoreRepository.RemoveStoreByID(ctx, tx, store.ID)
		if err != nil {
			errChan <- err
			return
		}

		err = iu.pgDBTx.CommitTx(tx)
		if err != nil {
			errChan <- err
			return
		}
		errChan <- nil
	}(errChan)

	err = <-errChan
	wg.Wait()

	if err != nil {
		return err
	}
	return nil
}

func (iu *integrationUsecase) CloseStoreRoutine(ctx context.Context) error {
	tx, err := iu.pgDBTx.BeginTx()
	if err != nil {
		return err
	}
	defer iu.pgDBTx.RollbackTx(tx)

	expires_time := time.Now().Add(-iu.config.StoreDuration)

	stores, err := iu.pgDBStoreRepository.GetAllExpiredStoresInSlice(ctx, tx, expires_time)
	if err != nil {
		return err
	}
	storeIds, err := iu.pgDBStoreRepository.GetAllIdsOfExpiredStores(ctx, tx, expires_time)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error)

	chunckedStores := iu.storeUsecase.ChunkStoresSlice(stores, 2)

	for _, stores := range chunckedStores {
		wg.Add(1)
		go func(stores [][][]string) {
			defer wg.Done()
			for _, store := range stores {
				storeInfo := store[0]
				store = store[1:]
				storeName, storeEmail, storeCreatedAtInstr, timezone := storeInfo[0], storeInfo[1], storeInfo[2], storeInfo[3]
				// str timestamp to int64 timestamp
				storeCreatedAtInInt64, _ := strconv.ParseInt(storeCreatedAtInstr, 10, 64)
				// timestamp to time
				storeCreatedAt := time.Unix(storeCreatedAtInInt64, 0)
				date, csvFileName, csvContent := iu.storeUsecase.GenerateCsvFileNameAndContent(storeCreatedAt, timezone, storeName, store)
				filePath, err := utils.GenerateCSV(
					ctx,
					iu.config.ContextTimeOut*4,
					csvFileName,
					csvContent,
				)
				if err != nil {
					return
				}
				emailSubject, emailContent := iu.storeUsecase.GenerateEmailContentOfCloseStore(storeName, date)
				_, _ = utils.SendEmail(
					ctx,
					iu.config.ContextTimeOut*4,
					iu.dialer,
					iu.config.FromEmail,
					emailSubject,
					emailContent,
					storeEmail,
					filePath,
				)
			}
		}(stores)
	}

	wg.Add(1)
	go func(errChan chan error, storeIds []string) {
		defer wg.Done()
		if len(storeIds) > 0 {
			// TODO: open it later...
			err := iu.pgDBStoreRepository.RemoveStoreByIDs(ctx, tx, storeIds)
			if err != nil {
				errChan <- err
				return
			}

			err = iu.pgDBTx.CommitTx(tx)
			if err != nil {
				errChan <- err
				return
			}
		}
		errChan <- nil
	}(errChan, storeIds)

	err = <-errChan
	wg.Wait()
	if err != nil {
		return err
	}

	return nil
}

func (iu *integrationUsecase) ForgetPassword(ctx context.Context, email string) (store domain.Store, err error) {
	store, err = iu.pgDBStoreRepository.GetStoreByEmail(ctx, email)
	if err != nil {
		return store, err
	}
	passwordToken, err := iu.storeUsecase.GenerateToken(
		ctx,
		store,
		domain.TokenTypes.PASSWORD,
		time.Now().Add(iu.config.PasswordTokenDuration),
	)
	if err != nil {
		return store, err
	}

	emailSubject, emailContent := iu.storeUsecase.GenerateEmailContentOfForgetPassword(passwordToken, store)
	_, err = utils.SendEmail(
		ctx,
		iu.config.ContextTimeOut*4,
		iu.dialer,
		iu.config.FromEmail,
		emailSubject,
		emailContent,
		email,
		"",
	)

	return store, err
}

func (iu *integrationUsecase) UpdatePassword(ctx context.Context, passwordToken string, newPassword string) (store domain.Store, err error) {
	tokenClaims, err := iu.storeUsecase.VerifyToken(
		ctx,
		passwordToken,
		domain.TokenTypes.PASSWORD,
		false,
	)
	if err != nil {
		return store, err
	}
	// TODO: delete all password tokens in tokens table

	store = domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}

	encryptedPassword, err := iu.storeUsecase.EncryptPassword(newPassword)
	if err != nil {
		return store, err
	}

	err = iu.pgDBStoreRepository.UpdateStore(ctx, &store, "password", encryptedPassword)
	if err != nil {
		return store, err
	}

	store.Password = encryptedPassword
	return store, nil
}

func (iu *integrationUsecase) GetStoreWithQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error) {
	store, err := iu.pgDBStoreRepository.GetStoreWithQueuesAndCustomersById(ctx, storeId)
	if err != nil {
		return store, err
	}
	if store.Queues == nil {
		return store, domain.ServerError40402
	}
	return store, err
}

func (iu *integrationUsecase) VerifyNormalToken(ctx context.Context, normalToken string) (tokenClaims domain.TokenClaims, err error) {
	encryptToken := strings.Split(normalToken, " ")
	if len(encryptToken) == 2 && strings.ToLower(encryptToken[0]) == "bearer" {
		tokenClaims, err = iu.storeUsecase.VerifyToken(
			ctx,
			encryptToken[1],
			domain.TokenTypes.NORMAL,
			true,
		)
		return tokenClaims, err
	}
	return tokenClaims, domain.ServerError40102
}

func (iu *integrationUsecase) VerifySessionToken(ctx context.Context, sessionToken string) (store domain.Store, err error) {
	tokenClaims, err := iu.storeUsecase.VerifyToken(
		ctx,
		sessionToken,
		domain.TokenTypes.SESSION,
		true, // TODO: change to false to RemoveTokenByToken
	)
	if err != nil {
		return store, err
	}
	store = domain.Store{
		ID:    tokenClaims.StoreID,
		Email: tokenClaims.Email,
		Name:  tokenClaims.Name,
	}
	return store, nil
}
