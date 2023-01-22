package pgDB

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBQueueRepository struct {
	db             PgDBInterface
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBQueueRepository(db PgDBInterface, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBQueueRepositoryInterface {
	return &pgDBQueueRepository{db, logger, contextTimeOut}
}

func (pqr *pgDBQueueRepository) CreateQueues(ctx context.Context, tx PgDBInterface, storeID int, queues []domain.Queue) error {
	ctx, cancel := context.WithTimeout(ctx, pqr.contextTimeOut)
	defer cancel()

	variableCounts := 1
	var query bytes.Buffer
	var queryRowParams []interface{}
	query.WriteString("INSERT INTO queues (name, store_id) VALUES ")
	for index, queue := range queues {
		query.WriteString("($")
		query.WriteString(strconv.Itoa(variableCounts))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 1))
		query.WriteString(")")
		variableCounts = variableCounts + 2
		queryRowParams = append(queryRowParams, queue.Name, storeID)
		if index != len(queues)-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(" RETURNING id,name")
	rows, err := tx.QueryContext(ctx, query.String(), queryRowParams...)
	queues = queues[:0] // clear queues slice

	for rows.Next() {
		queue := domain.Queue{}
		err = rows.Scan(&queue.ID, &queue.Name)
		if err != nil {
			pqr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		queue.StoreID = storeID
		queues = append(queues, queue)

	}
	defer rows.Close()

	return nil
}
