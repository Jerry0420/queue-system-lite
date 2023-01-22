package pgDB_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"github.com/stretchr/testify/assert"
)

func setUpSessionTest(t *testing.T) (pgDBSessionRepository pgDB.PgDBSessionRepositoryInterface, db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}

	logger := logging.NewLogger([]string{}, true)
	pgDBSessionRepository = pgDB.NewPgDBSessionRepository(db, logger, time.Duration(2*time.Second))
	return pgDBSessionRepository, db, mock
}

func TestGetSessionById(t *testing.T) {
	pgDBSessionRepository, _, mock := setUpSessionTest(t)
	mockStoreSession := domain.StoreSession{
		ID:                 "im_session_id",
		StoreId:            1,
		StoreSessionState: domain.StoreSessionState.NORMAL,
	}
	
	// store_id, session_state
	rows := sqlmock.NewRows([]string{"id", "state"}).
		AddRow(mockStoreSession.StoreId, mockStoreSession.StoreSessionState)

	query := `SELECT stores.id, store_sessions.state 
				FROM store_sessions
				INNER JOIN stores ON stores.id = store_sessions.store_id WHERE store_sessions.id=$1`
	mock.ExpectQuery(query).WithArgs(mockStoreSession.ID).WillReturnRows(rows)

	session, err := pgDBSessionRepository.GetSessionById(context.TODO(), mockStoreSession.ID)
	assert.NoError(t, err)
	assert.Equal(t, mockStoreSession.StoreId, session.StoreId)
	assert.Equal(t, mockStoreSession.ID, session.ID)
}
