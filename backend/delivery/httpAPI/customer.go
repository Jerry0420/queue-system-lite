package httpAPI

import (
	"net/http"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *HttpAPIDelivery) CustomersCreate(w http.ResponseWriter, r *http.Request) {
	session, customers, err := validator.CustomerCreate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.integrationUsecase.CreateCustomers(
		r.Context(),
		&session,
		domain.StoreSessionState.SCANNED,
		domain.StoreSessionState.USED,
		customers,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	go had.broker.Publish(
		had.storeUsecase.TopicNameOfUpdateCustomer(session.StoreId),
		map[string]interface{}{},
	)

	presenter.JsonResponseOK(w, customers)
}

// for used in store...
func (had *HttpAPIDelivery) CustomerUpdate(w http.ResponseWriter, r *http.Request) {
	storeId, oldCustomerState, newCustomerState, customer, err := validator.CustomerUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.customerUsecase.UpdateCustomer(r.Context(), oldCustomerState, newCustomerState, &customer)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	go had.broker.Publish(
		had.storeUsecase.TopicNameOfUpdateCustomer(storeId),
		map[string]interface{}{},
	)

	presenter.JsonResponseOK(w, customer)
}
