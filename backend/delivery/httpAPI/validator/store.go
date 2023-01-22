package validator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
)

func StoreOpen(r *http.Request) (store domain.Store, queues []domain.Queue, err error) {
	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return store, queues, domain.ServerError40001
	}
	name, ok := jsonBody["name"].(string)
	if !ok || name == "" {
		return store, queues, domain.ServerError40001
	}
	email, ok := jsonBody["email"].(string)
	if !ok || email == "" {
		return store, queues, domain.ServerError40001
	}
	password, ok := jsonBody["password"].(string)
	if !ok || password == "" {
		return store, queues, domain.ServerError40001
	}
	timezone, ok := jsonBody["timezone"].(string)
	if !ok || timezone == "" {
		return store, queues, domain.ServerError40001
	}
	store = domain.Store{Name: name, Email: email, Password: password, Timezone: timezone}
	
	queueNames, ok := jsonBody["queue_names"].([]interface{})
	if !ok || len(queueNames) <= 0 {
		return store, queues, domain.ServerError40001
	}
	for _, value := range queueNames {
		value, ok := value.(string)
		if !ok || value == "" {
			return store, queues, domain.ServerError40001
		}
		queues = append(queues, domain.Queue{Name: value})
	}
	return store, queues, nil
}

func StoreSignin(r *http.Request) (domain.Store, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return domain.Store{}, domain.ServerError40001
	}
	email, ok := jsonBody["email"].(string)
	if !ok || email == "" {
		return domain.Store{}, domain.ServerError40001
	}
	password, ok := jsonBody["password"].(string)
	if !ok || password == "" {
		return domain.Store{}, domain.ServerError40001
	}
	store := domain.Store{Email: email, Password: password}
	return store, nil
}

func StoreTokenRefresh(r *http.Request) (*http.Cookie, error) {
	encryptedRefreshToken, err := r.Cookie(domain.TokenTypes.REFRESH)
	if err != nil || len(encryptedRefreshToken.Value) == 0 {
		return nil, domain.ServerError40102
	}
	return encryptedRefreshToken, nil
}

func StoreClose(r *http.Request) (domain.TokenClaims, error) {
	tokenClaims := r.Context().Value(domain.TokenTypes.NORMAL).(domain.TokenClaims)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil || id != tokenClaims.StoreID {
		return domain.TokenClaims{}, domain.ServerError40004
	}
	return tokenClaims, nil
}

func StorePasswordForgot(r *http.Request) (domain.Store, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return domain.Store{}, domain.ServerError40001
	}
	email, ok := jsonBody["email"].(string)
	if !ok || email == "" {
		return domain.Store{}, domain.ServerError40001
	}
	store := domain.Store{Email: email}
	return store, nil
}

func StorePasswordUpdate(r *http.Request) (map[string]string, int, error) {
	var jsonBody map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return map[string]string{}, -1, domain.ServerError40001
	}
	passwordToken, ok := jsonBody["password_token"].(string)
	if !ok || passwordToken == "" {
		return map[string]string{}, -1, domain.ServerError40001
	}
	password, ok := jsonBody["password"].(string)
	if !ok || password == "" {
		return map[string]string{}, -1, domain.ServerError40001
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return map[string]string{}, -1, domain.ServerError40001
	}
	
	body := map[string]string{"password_token": passwordToken, "password": password} 
	return body, id, nil
}

func StoreInfoGet(r *http.Request) (storeId int, err error){
	vars := mux.Vars(r)
	storeId, err = strconv.Atoi(vars["id"])
	if err != nil {
		return storeId, domain.ServerError40001
	}
	return storeId, nil
}

func StoreDescriptionUpdate(r *http.Request) (store domain.Store, err error){
	tokenClaims := r.Context().Value(domain.TokenTypes.NORMAL).(domain.TokenClaims)
	
	vars := mux.Vars(r)
	storeId, err := strconv.Atoi(vars["id"])
	if err != nil || storeId != tokenClaims.StoreID {
		return store, domain.ServerError40004
	}

	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return store, domain.ServerError40001
	}
	description, ok := jsonBody["description"].(string)
	if !ok || description == "" {
		return store, domain.ServerError40001
	}

	store = domain.Store{
		ID: tokenClaims.StoreID,
		Email: tokenClaims.Email,
		Name: tokenClaims.Name,
		Description: description,
	}

	return store, nil
}