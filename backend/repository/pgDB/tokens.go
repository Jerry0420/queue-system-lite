package pgDB

import (
	"context"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
)

type pgDBTokenRepository struct {
	db             PgDBInterface
	logger         logging.LoggerTool
	contextTimeOut time.Duration
}

func NewPgDBTokenRepository(db PgDBInterface, logger logging.LoggerTool, contextTimeOut time.Duration) PgDBTokenRepositoryInterface {
	return &pgDBTokenRepository{db, logger, contextTimeOut}
}

func (pskr *pgDBTokenRepository) CreateToken(ctx context.Context, token *domain.Token) error {
	ctx, cancel := context.WithTimeout(ctx, pskr.contextTimeOut)
	defer cancel()

	query := `INSERT INTO tokens (store_id, token, type) VALUES ($1, $2, $3)`
	result, err := pskr.db.ExecContext(ctx, query, token.StoreId, token.Token, token.TokenType)
	if err != nil {
		pskr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	num, err := result.RowsAffected()
	if err != nil {
		pskr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	if num == 0 {
		return domain.ServerError40402 // foreignkey error
	}
	return nil
}

func (pskr *pgDBTokenRepository) RemoveTokenByToken(ctx context.Context, token string, tokenType string) error {
	ctx, cancel := context.WithTimeout(ctx, pskr.contextTimeOut)
	defer cancel()

	query := `DELETE FROM tokens WHERE token=$1 AND type=$2`
	result, err := pskr.db.ExecContext(ctx, query, token, tokenType)
	if err != nil {
		pskr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	num, err := result.RowsAffected()
	if err != nil {
		pskr.logger.ERRORf("error %v", err)
		return domain.ServerError50002
	}
	if num == 0 {
		return domain.ServerError40403
	}
	return nil
}
