package usecase

import (
	"context"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
)

type customerUsecase struct {
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface
	logger                 logging.LoggerTool
}

func NewCustomerUsecase(
	pgDBCustomerRepository pgDB.PgDBCustomerRepositoryInterface,
	logger logging.LoggerTool,
) CustomerUseCaseInterface {
	return &customerUsecase{pgDBCustomerRepository, logger}
}

func (cu *customerUsecase) UpdateCustomer(ctx context.Context, oldState string, newState string, customer *domain.Customer) error {
	err := cu.pgDBCustomerRepository.UpdateCustomer(ctx, oldState, newState, customer)
	return err
}
