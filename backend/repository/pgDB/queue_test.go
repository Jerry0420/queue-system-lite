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

func setUpQueueTest(t *testing.T) (pgDBQueueRepository pgDB.PgDBQueueRepositoryInterface, db *sql.DB, mock sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("error sqlmock new %v", err)
	}

	logger := logging.NewLogger([]string{}, true)
	pgDBQueueRepository = pgDB.NewPgDBQueueRepository(db, logger, time.Duration(2*time.Second))
	return pgDBQueueRepository, db, mock
}

func TestCreateQueues(t *testing.T) {
	pgDBQueueRepository, db, mock := setUpQueueTest(t)
	storeId := 1
	mockQueues := []domain.Queue{
		{
			Name:    "queue1",
			StoreID: storeId,
		},
		{
			Name:    "queue2",
			StoreID: storeId,
		},
	}

	query := `INSERT INTO queues (name, store_id) VALUES ($1, $2), ($3, $4) RETURNING id,name`
	mock.ExpectQuery(query).
		WithArgs(mockQueues[0].Name, mockQueues[0].StoreID, mockQueues[1].Name, mockQueues[1].StoreID).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(7, mockQueues[0].Name).AddRow(8, mockQueues[1].Name))

	err := pgDBQueueRepository.CreateQueues(context.TODO(), db, storeId, mockQueues)
	assert.NoError(t, err)
	assert.Equal(t, 7, mockQueues[0].ID)
	assert.Equal(t, "queue1", mockQueues[0].Name)
	assert.Equal(t, 8, mockQueues[1].ID)
	assert.Equal(t, "queue2", mockQueues[1].Name)
}
