package validator

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
)

func SessionCreate(r *http.Request) (sessionToken string, err error) {
	sessionToken = r.URL.Query().Get("session_token") // Because, it's a GET method.
	if sessionToken == "" {
		return sessionToken, domain.ServerError40001
	}
	return sessionToken, nil
}

func SessionScanned(r *http.Request) (session domain.StoreSession, err error) {
	session = r.Context().Value(domain.StoreSessionString).(domain.StoreSession)

	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return session, domain.ServerError40001
	}
	storeId, ok := jsonBody["store_id"].(float64)
	if !ok {
		return session, domain.ServerError40001
	}
	if int(storeId) != session.StoreId {
		return session, domain.ServerError40004
	}

	vars := mux.Vars(r)
	sessionId, ok := vars["id"] //sessionId
	if !ok || sessionId == "" || sessionId != session.ID {
		return session, domain.ServerError40004
	}

	return session, nil
}
