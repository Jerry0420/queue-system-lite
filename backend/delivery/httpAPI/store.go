package httpAPI

import (
	"fmt"
	"net/http"
	"time"

	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/presenter"
	"github.com/jerry0420/queue-system/backend/delivery/httpAPI/validator"
	"github.com/jerry0420/queue-system/backend/domain"
)

func (had *HttpAPIDelivery) OpenStore(w http.ResponseWriter, r *http.Request) {
	store, queues, err := validator.StoreOpen(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.storeUsecase.VerifyPasswordLength(store.Password)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.storeUsecase.VerifyTimeZoneString(store.Timezone)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.integrationUsecase.CreateStore(r.Context(), &store, queues)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreWithQueuesForResponse(store, queues))
}

func (had *HttpAPIDelivery) SigninStore(w http.ResponseWriter, r *http.Request) {
	incomingStore, err := validator.StoreSignin(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, token, refreshTokenExpiresAt, err := had.integrationUsecase.SigninStore(
		r.Context(),
		incomingStore.Email,
		incomingStore.Password,
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	cookieSecure := false
	
	if had.config.IsProdEnv == true {
		cookieSecure = true
	}
	cookieMaxAge := int(refreshTokenExpiresAt.Sub(time.Now()).Seconds())
	
	cookie := http.Cookie{
		Name:     domain.TokenTypes.REFRESH,
		Value:    token,
		Secure:   cookieSecure,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Path:     V_1("/stores/token"),
		MaxAge:   cookieMaxAge,
	}
	http.SetCookie(w, &cookie)

	// use cookie: refreshable to check if refresh token exist in the browser. 
	cookie = http.Cookie{
		Name:     domain.TokenTypes.REFRESH + "able",
		Value:    "true",
		Secure:   cookieSecure,
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   cookieMaxAge,
	}
	http.SetCookie(w, &cookie)

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *HttpAPIDelivery) RefreshToken(w http.ResponseWriter, r *http.Request) {
	encryptedRefreshToken, err := validator.StoreTokenRefresh(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, normalToken, sessionToken, tokenExpiresAt, err := had.integrationUsecase.RefreshToken(r.Context(), encryptedRefreshToken.Value)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreToken(store, normalToken, tokenExpiresAt, sessionToken))
}

func (had *HttpAPIDelivery) CloseStore(w http.ResponseWriter, r *http.Request) {
	tokenClaims, err := validator.StoreClose(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store := domain.Store{
		ID:        tokenClaims.StoreID,
		Email:     tokenClaims.Email,
		Name:      tokenClaims.Name,
		CreatedAt: time.Unix(tokenClaims.StoreCreatedAt, 0),
	}
	err = had.integrationUsecase.CloseStore(r.Context(), store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *HttpAPIDelivery) CloseStorerRoutine(w http.ResponseWriter, r *http.Request) {
	err := had.integrationUsecase.CloseStoreRoutine(r.Context())
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, nil)
}

func (had *HttpAPIDelivery) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	store, err := validator.StorePasswordForgot(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err = had.integrationUsecase.ForgetPassword(r.Context(), store.Email)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *HttpAPIDelivery) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	jsonBody, _, err := validator.StorePasswordUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.storeUsecase.VerifyPasswordLength(jsonBody["password"])
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err := had.integrationUsecase.UpdatePassword(
		r.Context(),
		jsonBody["password_token"],
		jsonBody["password"],
	)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}

func (had *HttpAPIDelivery) GetStoreInfoWithSSE(w http.ResponseWriter, r *http.Request) {
	storeId, err := validator.StoreInfoGet(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	flusher, ok := w.(http.Flusher)
	if !ok {
		presenter.JsonResponse(w, nil, domain.ServerError50003)
		return
	}
	consumerChan := had.broker.Subscribe(had.storeUsecase.TopicNameOfUpdateCustomer(storeId))
	defer had.broker.UnsubscribeConsumer(had.storeUsecase.TopicNameOfUpdateCustomer(storeId), consumerChan)

	store, err := had.integrationUsecase.GetStoreWithQueuesAndCustomersById(r.Context(), storeId)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	fmt.Fprintf(w, "data: %v\n\n", presenter.StoreGetForSSE(store))
	flusher.Flush()

	for {
		select {
		case <-consumerChan:
			store, err := had.integrationUsecase.GetStoreWithQueuesAndCustomersById(r.Context(), storeId)
			if err != nil {
				presenter.JsonResponse(w, nil, err)
				return
			}
			fmt.Fprintf(w, "data: %v\n\n", presenter.StoreGetForSSE(store))
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func (had *HttpAPIDelivery) GetStoreInfo(w http.ResponseWriter, r *http.Request) {
	storeId, err := validator.StoreInfoGet(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	store, err := had.integrationUsecase.GetStoreWithQueuesAndCustomersById(r.Context(), storeId)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}
	presenter.JsonResponseOK(w, presenter.StoreGet(store))
}

func (had *HttpAPIDelivery) UpdateStoreDescription(w http.ResponseWriter, r *http.Request) {
	store, err := validator.StoreDescriptionUpdate(r)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	err = had.storeUsecase.UpdateStoreDescription(r.Context(), store.Description, &store)
	if err != nil {
		presenter.JsonResponse(w, nil, err)
		return
	}

	go had.broker.Publish(
		had.storeUsecase.TopicNameOfUpdateCustomer(store.ID),
		map[string]interface{}{},
	)

	presenter.JsonResponseOK(w, presenter.StoreForResponse(store))
}
