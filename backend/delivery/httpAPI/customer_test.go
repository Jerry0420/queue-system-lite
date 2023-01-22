package httpAPI_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCustomerUpdate(t *testing.T) {
	customerUseCase, _, storeUseCase, _, httpAPIDelivery, router, broker := setUpHttpAPITest()
	defer broker.CloseAll()

	originalCustomerState := domain.CustomerState.WAITING
	expectedCustomerState := domain.CustomerState.PROCESSING

	mockCustomer := domain.Customer{
		ID: 1,
		QueueID: 1,
		State:  expectedCustomerState,
	}

	mockTokenClaims := domain.TokenClaims{
		StoreID: 1,
		Email:   "storeemail",
		Name:    "storename",
	}
	customerUseCase.On("UpdateCustomer", mock.Anything, originalCustomerState, expectedCustomerState, &mockCustomer).Return(nil).Once()
	storeUseCase.On("TopicNameOfUpdateCustomer", mockTokenClaims.StoreID).Return("im_topic").Once()

	router.HandleFunc(
		httpAPI.V_1("/customers/{id:[0-9]+}"),
		httpAPIDelivery.CustomerUpdate,
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	ctx := context.WithValue(context.Background(), domain.TokenTypes.NORMAL, mockTokenClaims)
	w := httptest.NewRecorder()
	params := map[string]interface{}{
		"store_id":            mockTokenClaims.StoreID,
		"queue_id":            mockCustomer.QueueID,
		"old_customer_state": originalCustomerState,
		"new_customer_state": expectedCustomerState,
	}
	jsonBody, _ := json.Marshal(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/api/v1/customers/"+strconv.Itoa(mockCustomer.ID), bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	router.ServeHTTP(w, req)

	var decodedResponse map[string]interface{}
	json.NewDecoder(w.Result().Body).Decode(&decodedResponse)
	assert.Equal(t, expectedCustomerState, decodedResponse["state"].(string))
}
