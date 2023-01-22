package domain

import "time"

type customerState struct{ WAITING, PROCESSING, DONE, DELETE string }

var CustomerState customerState = customerState{WAITING: "waiting", PROCESSING: "processing", DONE: "done", DELETE: "delete"}

type Customer struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	QueueID   int       `json:"-"`
	State     string    `json:"state"`
	CreatedAt time.Time `json:"created_at"`
}
