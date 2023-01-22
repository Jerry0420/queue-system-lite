package validator

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jerry0420/queue-system/backend/domain"
)

func CustomerCreate(r *http.Request) (session domain.StoreSession, customers []domain.Customer, err error) {
	session = r.Context().Value(domain.StoreSessionString).(domain.StoreSession)

	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return session, customers, domain.ServerError40001
	}

	storeIdFloat64, ok := jsonBody["store_id"].(float64)
	if !ok {
		return session, customers, domain.ServerError40001
	}

	customersInfo, ok := jsonBody["customers"].([]interface{})
	if !ok || len(customersInfo) <= 0 {
		return session, customers, domain.ServerError40001
	}
// exceed 5 customers at a req is forbidden.
	if len(customersInfo) > 5 {
		return session, customers, domain.ServerError40005
	}

	for _, customerInfo := range customersInfo {
		customerInfo, ok := customerInfo.(map[string]interface{})
		if !ok {
			return session, customers, domain.ServerError40001
		}
		name, ok := customerInfo["name"].(string)
		if !ok || name == "" {
			return session, customers, domain.ServerError40001
		}
		phone, ok := customerInfo["phone"].(string)
		
		queueId, ok := customerInfo["queue_id"].(float64)
		if !ok {
			return session, customers, domain.ServerError40001
		}
		customers = append(
			customers, 
			domain.Customer{Name: name, Phone: phone, QueueID: int(queueId), State: domain.CustomerState.WAITING},
		)
	}

	if int(storeIdFloat64) != session.StoreId {
		return session, customers, domain.ServerError40004
	}

	return session, customers, nil
}

func CustomerUpdate(r *http.Request) (storeId int, oldCustomerState string, newCustomerState string, customer domain.Customer, err error) {
	tokenClaims := r.Context().Value(domain.TokenTypes.NORMAL).(domain.TokenClaims)

	vars := mux.Vars(r)
	customerId, err := strconv.Atoi(vars["id"])
	if err != nil {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40001
	}

	var jsonBody map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&jsonBody)
	if err != nil {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40001
	}

	storeIdFloat64, ok := jsonBody["store_id"].(float64)
	if !ok {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40001
	}

	if int(storeIdFloat64) != tokenClaims.StoreID {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40004
	}

	queueIdFloat64, ok := jsonBody["queue_id"].(float64)
	if !ok {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40001
	}

	oldCustomerState, ok = jsonBody["old_customer_state"].(string)
	if !ok || oldCustomerState == "" {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40001
	}

	newCustomerState, ok = jsonBody["new_customer_state"].(string)
	if !ok || newCustomerState == "" {
		return storeId, oldCustomerState, newCustomerState, customer, domain.ServerError40001
	}

	customer = domain.Customer{
		ID: customerId,
		QueueID: int(queueIdFloat64),
		State: newCustomerState,
	}

	return int(storeIdFloat64), oldCustomerState, newCustomerState, customer, nil
}
