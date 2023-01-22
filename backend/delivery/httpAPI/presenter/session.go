package presenter

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/jerry0420/queue-system/backend/domain"
)

func SessionCreate(domain string, session domain.StoreSession) string {
	var sessionJson []byte
	var sessionMap map[string]interface{}
	sessionJson, _ = json.Marshal(session)
	json.Unmarshal(sessionJson, &sessionMap)
	sessionMap["scanned_url"] = fmt.Sprintf("%s/#/stores/%d/sessions/%s", domain, session.StoreId, session.ID)
	var flushedData bytes.Buffer
	json.NewEncoder(&flushedData).Encode(sessionMap)
	return flushedData.String()
}
