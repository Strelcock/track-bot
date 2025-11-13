package kservice

import "time"

type Status struct {
	Event     string    `json:"event"`
	Barcode   string    `json:"barcode"`
	Status    string    `json:"status"`
	EventDate time.Time `json:"eventDate"`
}

type StatusTo struct {
	Status *Status
	To     []int64
}
