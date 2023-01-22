package pgDB

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBSessionRepository struct {
	db             PgDBInterface
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBSessionRepository(db PgDBInterface, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBSessionRepositoryInterface {
	return &pgDBSessionRepository{db, logger, contextTimeOut}
}

func (psr *pgDBSessionRepository) CreateSession(ctx context.Context, store domain.Store) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{StoreId: store.ID, StoreSessionState: domain.StoreSessionState.NORMAL}

	query := `INSERT INTO store_sessions (store_id) VALUES ($1) RETURNING id`
	row := psr.db.QueryRowContext(ctx, query, store.ID)
	err := row.Scan(&session.ID)
	if err != nil {
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError40903
	}
	return session, nil
}

func (psr *pgDBSessionRepository) UpdateSessionState(ctx context.Context, tx PgDBInterface, session *domain.StoreSession, oldState string, newState string) error {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	sessionStateInDb := ""
	var err error
	var row *sql.Row

	query := `SELECT state FROM store_sessions WHERE id=$1`
	if tx == nil {
		row = psr.db.QueryRowContext(ctx, query, session.ID)
	} else {
		row = tx.QueryRowContext(ctx, query, session.ID)
	}
	err = row.Scan(&sessionStateInDb)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		psr.logger.ERRORf("error %v", err)
		return domain.ServerError40404
	case err != nil:
		psr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}

	if sessionStateInDb == oldState {
		query = `UPDATE store_sessions SET state=$1 WHERE id=$2 and state=$3`
		var result sql.Result
		if tx == nil {
			result, err = psr.db.ExecContext(ctx, query, newState, session.ID, oldState)
		} else {
			result, err = tx.ExecContext(ctx, query, newState, session.ID, oldState)
		}
		if err != nil {
			psr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		num, err := result.RowsAffected()
		if err != nil {
			psr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		if num == 0 {
			return domain.ServerError40404
		}
		return nil
	} else {
		switch sessionStateInDb {
			case domain.StoreSessionState.SCANNED: 
				return domain.ServerError40007
			case domain.StoreSessionState.USED: 
				return domain.ServerError40008
		}
	}

	return nil
}

func (psr *pgDBSessionRepository) GetSessionById(ctx context.Context, sessionId string) (domain.StoreSession, error) {
	ctx, cancel := context.WithTimeout(ctx, psr.contextTimeOut)
	defer cancel()

	session := domain.StoreSession{}
	store := domain.Store{}

	query := `SELECT stores.id, store_sessions.state 
				FROM store_sessions
				INNER JOIN stores ON stores.id = store_sessions.store_id WHERE store_sessions.id=$1`
	row := psr.db.QueryRowContext(ctx, query, sessionId)
	err := row.Scan(&store.ID, &session.StoreSessionState)
	switch {
	case errors.Is(err, sql.ErrNoRows):
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError40404
	case err != nil:
		psr.logger.ERRORf("error %v", err)
		return session, domain.ServerError50002
	}
	session.ID = sessionId
	session.StoreId = store.ID
	return session, nil
}
