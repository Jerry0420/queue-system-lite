package domain

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Store struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Password    string    `json:"-"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	Timezone    string    `json:"-"`
}

type TokenClaims struct {
	StoreID        int    `json:"store_id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	StoreCreatedAt int64  `json:"store_created_at"` // not time.Time
	jwt.StandardClaims
}

type StoreWithQueues struct {
	ID          int                   `json:"id"`
	Email       string                `json:"email"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	CreatedAt   time.Time             `json:"created_at"`
	Queues      []*QueueWithCustomers `json:"queues"`
}
