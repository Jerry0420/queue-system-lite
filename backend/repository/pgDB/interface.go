package pgDB

import (
	"context"
	"database/sql"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
)

// interface for sql.DB and sql.Tx
type PgDBInterface interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
	Prepare(query string) (*sql.Stmt, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	QueryRow(query string, args ...interface{}) *sql.Row
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type PgDBTxInterface interface {
	BeginTx() (tx PgDBInterface, err error)
	RollbackTx(pgDbTx PgDBInterface)
	CommitTx(pgDbTx PgDBInterface) error
}

type PgDBStoreRepositoryInterface interface {
	GetStoreByEmail(ctx context.Context, email string) (domain.Store, error)
	GetStoreWithQueuesAndCustomersById(ctx context.Context, storeId int) (domain.StoreWithQueues, error)
	CreateStore(ctx context.Context, tx PgDBInterface, store *domain.Store) error
	UpdateStore(ctx context.Context, store *domain.Store, fieldName string, newFieldValue string) error
	RemoveStoreByID(ctx context.Context, tx PgDBInterface, id int) error
	RemoveStoreByIDs(ctx context.Context, tx PgDBInterface, storeIds []string) error
	GetAllIdsOfExpiredStores(ctx context.Context, tx PgDBInterface, expiresTime time.Time) (storesIds []string, err error)
	GetAllExpiredStoresInSlice(ctx context.Context, tx PgDBInterface, expiresTime time.Time) (stores [][][]string, err error)
}

type PgDBTokenRepositoryInterface interface {
	CreateToken(ctx context.Context, token *domain.Token) error
	RemoveTokenByToken(ctx context.Context, token string, tokenType string) error
}

type PgDBQueueRepositoryInterface interface {
	CreateQueues(ctx context.Context, tx PgDBInterface, storeID int, queues []domain.Queue) error
}

type PgDBCustomerRepositoryInterface interface {
	CreateCustomers(ctx context.Context, tx PgDBInterface, customers []domain.Customer) error
	UpdateCustomer(ctx context.Context, oldState string, newState string, customer *domain.Customer) error
	GetCustomersWithQueuesByStore(ctx context.Context, tx PgDBInterface, store *domain.Store) (customers [][]string, err error)
}

type PgDBSessionRepositoryInterface interface {
	CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error)
	UpdateSessionState(ctx context.Context, tx PgDBInterface, session *domain.StoreSession, oldState string, newState string) error
	GetSessionById(ctx context.Context, sessionId string) (domain.StoreSession, error)
}
