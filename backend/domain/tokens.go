package domain

type tokenTypes struct{ NORMAL, PASSWORD, REFRESH, SESSION string }

var TokenTypes tokenTypes = tokenTypes{NORMAL: "normal", PASSWORD: "password", REFRESH: "refresh", SESSION: "session"}

type Token struct {
	StoreId   int    `json:"-"`
	Token     string `json:"token"`
	TokenType string `json:"type"`
}
