package usecase

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/jerry0420/queue-system/backend/logging"
	"github.com/jerry0420/queue-system/backend/repository/pgDB"
	"golang.org/x/crypto/bcrypt"
)

type StoreUsecaseConfig struct {
	Domain       string
	TokenSignKey string
}

type storeUsecase struct {
	pgDBStoreRepository pgDB.PgDBStoreRepositoryInterface
	pgDBTokenRepository pgDB.PgDBTokenRepositoryInterface
	logger              logging.LoggerTool
	config              StoreUsecaseConfig
}

func NewStoreUsecase(
	pgDBStoreRepository pgDB.PgDBStoreRepositoryInterface,
	pgDBTokenRepository pgDB.PgDBTokenRepositoryInterface,
	logger logging.LoggerTool,
	config StoreUsecaseConfig,
) StoreUseCaseInterface {
	return &storeUsecase{pgDBStoreRepository, pgDBTokenRepository, logger, config}
}

func (su *storeUsecase) ChunkStoresSlice(items [][][]string, chunkSize int) (chunks [][][][]string) {
	for chunkSize < len(items) {
		chunks = append(chunks, items[0:chunkSize])
		items = items[chunkSize:]
	}
	return append(chunks, items)
}

func (su *storeUsecase) UpdateStoreDescription(ctx context.Context, newDescription string, store *domain.Store) error {
	err := su.pgDBStoreRepository.UpdateStore(ctx, store, "description", newDescription)
	if err != nil {
		return err
	}
	return nil
}

func (su *storeUsecase) VerifyPasswordLength(password string) error {
	decodedPassword, err := base64.StdEncoding.DecodeString(password)
	if err != nil {
		su.logger.ERRORf("%v", err)
		return domain.ServerError50001
	}
	rawPassword := string(decodedPassword)
	// length of password must between 8 and 15.
	if len(rawPassword) < 8 || len(rawPassword) > 15 {
		return domain.ServerError40002
	}
	return nil
}

func (su *storeUsecase) VerifyTimeZoneString(inputTimezone string) error {
	_, err := time.LoadLocation(inputTimezone)
	if err != nil {
		su.logger.ERRORf("%v", err)
		return domain.ServerError40006
	}
	return nil
}

func (su *storeUsecase) EncryptPassword(password string) (string, error) {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		su.logger.ERRORf("%v", err)
		return "", domain.ServerError50001
	}
	return string(cryptedPassword), nil
}

func (su *storeUsecase) ValidatePassword(passwordInDb string, incomingPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(passwordInDb), []byte(incomingPassword))
	switch {
	case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
		su.logger.ERRORf("%v", err)
		return domain.ServerError40003
	case err != nil:
		su.logger.ERRORf("%v", err)
		return domain.ServerError50001
	}
	return nil
}

func (su *storeUsecase) GenerateToken(ctx context.Context, store domain.Store, tokenType string, expireTime time.Time) (encryptToken string, err error) {
	// randomUUID := uuid.New().String()
	// saltBytes, err := bcrypt.GenerateFromPassword([]byte(randomUUID), bcrypt.DefaultCost)
	// if err != nil {
	// 	su.logger.ERRORf("%v", err)
	// 	return "", domain.ServerError50001
	// }

	claims := domain.TokenClaims{
		store.ID,
		store.Email,
		store.Name,
		store.CreatedAt.Unix(),
		jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: expireTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	encryptToken, err = token.SignedString([]byte(su.config.TokenSignKey))
	if err != nil {
		su.logger.ERRORf("%v", err)
		return encryptToken, domain.ServerError50001
	}
	err = su.pgDBTokenRepository.CreateToken(
		ctx,
		&domain.Token{StoreId: store.ID, Token: encryptToken, TokenType: tokenType},
	)
	if err != nil {
		return "", err
	}
	return encryptToken, err
}

func (su *storeUsecase) VerifyToken(ctx context.Context, encryptToken string, tokenType string, withTokenPreserved bool) (tokenClaims domain.TokenClaims, err error) {
	tokenClaims = domain.TokenClaims{}
	token, err := jwt.ParseWithClaims(encryptToken, &tokenClaims, func(token *jwt.Token) (interface{}, error) {
		if !withTokenPreserved { // == false
			su.pgDBTokenRepository.RemoveTokenByToken(ctx, encryptToken, tokenType)
		}
		return []byte(su.config.TokenSignKey), nil
	})
	if err != nil {
		su.logger.ERRORf("%v", err)
		if err.(*jwt.ValidationError).Errors == jwt.ValidationErrorExpired {
			return tokenClaims, domain.ServerError40103
		}
		if serverError, ok := err.(*jwt.ValidationError).Inner.(*domain.ServerError); ok {
			return tokenClaims, serverError
		}
		return tokenClaims, domain.ServerError40101
	}

	if !token.Valid {
		su.logger.ERRORf("unvalid token")
		return tokenClaims, domain.ServerError40101
	}

	return tokenClaims, nil
}

func (su *storeUsecase) GenerateEmailContentOfForgetPassword(passwordToken string, store domain.Store) (subject string, content string) {
	// TODO: update email content to html format.
	resetPasswordUrl := fmt.Sprintf("%s/#/stores/%d/password/update?password_token=%s", su.config.Domain, store.ID, passwordToken)
	return "Queue-System Reset Password", fmt.Sprintf("Hello, %s, please click %s", store.Name, resetPasswordUrl)
}

func (su *storeUsecase) GenerateEmailContentOfCloseStore(storeName string, storeCreatedAt string) (subject string, content string) {
	// TODO: update email content to html format.
	subject = fmt.Sprintf("Queue-System: Result of %s (%s)", storeName, storeCreatedAt)
	content = fmt.Sprintf("Hello %s, The attached file is the result of %s.\n\nThank you", storeName, storeCreatedAt)
	return subject, content
}

func (su *storeUsecase) GenerateCsvFileNameAndContent(storeCreatedAt time.Time, storeTimezone string, storeName string, content [][]string) (date string, csvFileName string, csvContent []byte) {
	timezone, _ := time.LoadLocation(storeTimezone)
	year, month, day := storeCreatedAt.UTC().In(timezone).Date()
	date = fmt.Sprintf("%d-%d-%d", year, month, day)
	csvFileName = fmt.Sprintf("%s-%s", date, storeName)
	csvContent, _ = json.Marshal(content)
	return date, csvFileName, csvContent
}

func (su *storeUsecase) TopicNameOfUpdateCustomer(storeId int) string {
	return fmt.Sprintf("updateCustomer.%d", storeId)
}
