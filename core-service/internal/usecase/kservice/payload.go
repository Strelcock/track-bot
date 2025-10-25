package kservice

import "time"

type Status struct {
	ID        string    `json:"id"`
	Status    string    `json:"status"`
	UpdatedAt time.Time `json:"updated_at"`
}

type StatusTo struct {
	Status Status
	To     int64
}
