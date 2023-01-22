package domain

type storeSessionState struct{ NORMAL, SCANNED, USED string }

var StoreSessionState storeSessionState = storeSessionState{NORMAL: "normal", SCANNED: "scanned", USED: "used"}

const StoreSessionString string = "session"

type StoreSession struct {
	ID                string `json:"id"`
	StoreId           int    `json:"-"`
	StoreSessionState string `json:"state"`
}
