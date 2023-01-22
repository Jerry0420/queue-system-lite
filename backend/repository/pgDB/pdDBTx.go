package pgDB

import (
	"database/sql"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBTx struct {
	db     PgDBInterface
	logger logging.LoggerTool
}

func NewPgDBTx(db PgDBInterface, logger logging.LoggerTool) PgDBTxInterface {
	return &pgDBTx{db, logger}
}

func (pt *pgDBTx) BeginTx() (pgDbTxHandle PgDBInterface, err error) {
	// without ctx
	tx, err := pt.db.(*sql.DB).Begin()
	if err != nil {
		pt.logger.ERRORf("begin tx error %v", err)
		return nil, domain.ServerError50002
	}
	return tx, nil
}

func (pt *pgDBTx) RollbackTx(pgDbTx PgDBInterface) {
	_ = pgDbTx.(*sql.Tx).Rollback()
}

func (pt *pgDBTx) CommitTx(pgDbTx PgDBInterface) error {
	err := pgDbTx.(*sql.Tx).Commit()
	if err != nil {
		pt.logger.ERRORf("commit error %v", err)
		return domain.ServerError50002
	}
	return nil
}
