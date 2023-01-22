package httpAPI_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI"
	"github.com/jerry0420/queue-system/backend/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestScannedSession(t *testing.T) {
	_, sessionUseCase, _, _, httpAPIDelivery, router, broker := setUpHttpAPITest()
	defer broker.CloseAll()
	expectedMockSessionState := domain.StoreSessionState.SCANNED
	mockSession := domain.StoreSession{
		ID:                 "im_session_id",
		StoreId:            1,
		StoreSessionState: domain.StoreSessionState.NORMAL,
	}

	sessionUseCase.On("UpdateSessionState", mock.Anything, &mockSession, mockSession.StoreSessionState, expectedMockSessionState).
		Return(nil).Run(func(args mock.Arguments) {
			session := args.Get(1).(*domain.StoreSession)
			newSessionState := args.Get(3).(string)
			session.StoreSessionState = newSessionState
		}).Once()
	sessionUseCase.On("TopicNameOfUpdateSession", mockSession.StoreId).Return("im_topic").Once()

	router.HandleFunc(
		httpAPI.V_1("/sessions/{id}"),
		httpAPIDelivery.ScannedSession,
	).Methods(http.MethodPut).Headers("Content-Type", "application/json")

	ctx := context.WithValue(context.Background(), domain.StoreSessionString, mockSession)
	w := httptest.NewRecorder()
	params := map[string]interface{}{
		"store_id": mockSession.StoreId,
	}
	jsonBody, _ := json.Marshal(params)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, "/api/v1/sessions/"+mockSession.ID, bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	assert.NoError(t, err)
	router.ServeHTTP(w, req)

	var decodedResponse map[string]interface{}
	json.NewDecoder(w.Result().Body).Decode(&decodedResponse)
	assert.Equal(t, expectedMockSessionState, decodedResponse["state"].(string))
}
