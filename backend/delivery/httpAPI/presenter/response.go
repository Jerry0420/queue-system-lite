package presenter

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/jerry0420/queue-system/backend/config"
	"github.com/jerry0420/queue-system/backend/domain"
)

func JsonResponseOK(w http.ResponseWriter, response interface{}) {
	JsonResponse(w, response, nil)
}

func JsonResponse(w http.ResponseWriter, response interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")

	if err != nil {
		serverError, ok := err.(*domain.ServerError)
		if !ok {
			// Because err must be ServerError type.
			// Using panic to prevent any hidden error when develping...
			panic("err must be ServerError")
		}
		errorResponse := map[string]interface{}{"error_code": serverError.Code}
		if config.ServerConfig.ENV() != config.EnvState.PROD {
			errorResponse["description"] = serverError.Description
		}

		w.WriteHeader(serverError.Code / 100) //http status code is the first two digits of code.
		w.Header().Set("Server-Code", strconv.Itoa(serverError.Code))
		json.NewEncoder(w).Encode(&errorResponse)
	} else {
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(&response)
	}
}
