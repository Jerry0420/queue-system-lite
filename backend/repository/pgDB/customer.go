package pgDB

import (
	"bytes"
	"context"
	"strconv"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBCustomerRepository struct {
	db             PgDBInterface
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBCustomerRepository(db PgDBInterface, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBCustomerRepositoryInterface {
	return &pgDBCustomerRepository{db, logger, contextTimeOut}
}

func (pcr *pgDBCustomerRepository) CreateCustomers(ctx context.Context, tx PgDBInterface, customers []domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, pcr.contextTimeOut)
	defer cancel()

	variableCounts := 1
	var query bytes.Buffer
	var queryRowParams []interface{}
	query.WriteString("INSERT INTO customers (name, phone, queue_id, state) VALUES ")
	for index, customer := range customers {
		query.WriteString("($")
		query.WriteString(strconv.Itoa(variableCounts))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 1))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 2))
		query.WriteString(", $")
		query.WriteString(strconv.Itoa(variableCounts + 3))
		query.WriteString(")")
		variableCounts = variableCounts + 4
		queryRowParams = append(queryRowParams, customer.Name, customer.Phone, customer.QueueID, customer.State)
		if index != len(customers)-1 {
			query.WriteString(", ")
		}
	}
	query.WriteString(" RETURNING id,name,phone,queue_id,created_at")

	rows, err := tx.QueryContext(ctx, query.String(), queryRowParams...)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	customers = customers[:0] // clear customers slice

	for rows.Next() {
		customer := domain.Customer{}
		err = rows.Scan(&customer.ID, &customer.Name, &customer.Phone, &customer.QueueID, &customer.CreatedAt)
		if err != nil {
			pcr.logger.ERRORf("error %v", err)
			return domain.ServerError50002
		}
		customer.State = domain.CustomerState.WAITING
		customers = append(customers, customer)

	}
	defer rows.Close()

	return nil
}

func (pcr *pgDBCustomerRepository) UpdateCustomer(ctx context.Context, oldState string, newState string, customer *domain.Customer) error {
	ctx, cancel := context.WithTimeout(ctx, pcr.contextTimeOut)
	defer cancel()

	query := `UPDATE customers SET state=$1 WHERE id=$2 and state=$3`
	result, err := pcr.db.ExecContext(ctx, query, newState, customer.ID, oldState)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	num, err := result.RowsAffected()
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	if num == 0 {
		return domain.ServerError40405
	}
	return nil
}

func (pcr *pgDBCustomerRepository) GetCustomersWithQueuesByStore(ctx context.Context, tx PgDBInterface, store *domain.Store) (customers [][]string, err error) {
	ctx, cancel := context.WithTimeout(ctx, pcr.contextTimeOut)
	defer cancel()

	customers = make([][]string, 0)

	query := `SELECT 
					stores.timezone AS timezone, 
					queues.name AS queue_name, 
					customers.name AS customer_name, customers.phone AS customer_phone,
					customers.state AS customer_state, customers.created_at AS customer_created_at
				FROM queues
				INNER JOIN customers ON queues.id = customers.queue_id
				INNER JOIN stores ON stores.id = queues.store_id
				WHERE queues.store_id=$1
				ORDER BY queues.id ASC, customers.id ASC FOR UPDATE` // row block

	rows, err := tx.QueryContext(ctx, query, store.ID)
	if err != nil {
		pcr.logger.ERRORf("error %v", err)
		return customers, domain.ServerError50002
	}

	for rows.Next() {
		var queue domain.Queue
		var customer domain.Customer

		err := rows.Scan(
			&store.Timezone,
			&queue.Name,
			&customer.Name, &customer.Phone, &customer.State, &customer.CreatedAt,
		)
		if err != nil {
			pcr.logger.ERRORf("error %v", err)
			return customers, domain.ServerError50002
		}

		timezone, _ := time.LoadLocation(store.Timezone)

		if len(customers) == 0 {
			customers = [][]string{
				[]string{
					"queue_name",
					"customer_name",
					"customer_phone",
					"customer_state",
					"customer_created_at",
				},
				[]string{
					queue.Name,
					customer.Name,
					customer.Phone,
					customer.State,
					customer.CreatedAt.UTC().In(timezone).String(),
				},
			}
		} else {
			customers = append(customers, []string{
				queue.Name,
				customer.Name,
				customer.Phone,
				customer.State,
				customer.CreatedAt.UTC().In(timezone).String(),
			})
		}
	}
	defer rows.Close()
	return customers, nil
}
