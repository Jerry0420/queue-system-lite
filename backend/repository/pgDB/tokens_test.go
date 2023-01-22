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

func setUpTokensTest(t *testing.T) (pgDBTokenRepository pgDB.PgDBTokenRepositoryInterface, db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}

	logger := logging.NewLogger([]string{}, true)
	pgDBTokenRepository = pgDB.NewPgDBTokenRepository(db, logger, time.Duration(2*time.Second))
	return pgDBTokenRepository, db, mock
}

func TestCreateToken(t *testing.T) {
	pgDBTokenRepository, _, mock := setUpTokensTest(t)
	mockToken := domain.Token{
		StoreId: 1,
		Token: "imtoken",
		TokenType: domain.TokenTypes.NORMAL,
	}

	query := `INSERT INTO tokens (store_id, token, type) VALUES ($1, $2, $3)`
	mock.ExpectExec(query).
		WithArgs(mockToken.StoreId, mockToken.Token, mockToken.TokenType).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := pgDBTokenRepository.CreateToken(context.TODO(), &mockToken)
	assert.NoError(t, err)
}
