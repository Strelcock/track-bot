package kservice

type Status struct {
	Event   string   `json:"event"`
	Barcode string   `json:"barcode"`
	Status  []string `json:"status"`
}

type StatusTo struct {
	Status *Status
	To     []int64
}
