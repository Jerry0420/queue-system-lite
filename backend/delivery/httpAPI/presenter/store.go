package presenter

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/jerry0420/queue-system/backend/domain"
)

func StoreWithQueuesForResponse(store domain.Store, queues []domain.Queue) map[string]interface{} {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	// created_at, _ := time.Parse(time.RFC3339, storeMap["created_at"].(string))
	// storeMap["created_at"] = created_at.Unix()

	var queuesJson []byte
	var queuesMap []map[string]interface{}
	queuesJson, _ = json.Marshal(queues)
	json.Unmarshal(queuesJson, &queuesMap)
	storeMap["queues"] = queuesMap
	return storeMap
}

func StoreForResponse(store domain.Store) map[string]interface{} {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	return storeMap
}

func StoreToken(store domain.Store, normalToken string, tokenExpiresAt time.Time, sessionToken string) map[string]interface{} {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	storeMap["token"] = normalToken
	storeMap["token_expires_at"] = tokenExpiresAt.Unix()
	storeMap["session_token"] = sessionToken
	return storeMap
}

func StoreGetForSSE(store domain.StoreWithQueues) string {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	var flushedData bytes.Buffer
	json.NewEncoder(&flushedData).Encode(storeMap)
	return flushedData.String()
}

func StoreGet(store domain.StoreWithQueues) map[string]interface{} {
	var storeJson []byte
	var storeMap map[string]interface{}
	storeJson, _ = json.Marshal(store)
	json.Unmarshal(storeJson, &storeMap)
	return storeMap
}
