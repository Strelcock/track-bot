package hook

type CarrierResponse struct {
	Code string `json:"code"`
}

type NewEvent struct {
	NewEvents EventData `json:"newEvents"`
}

type EventData struct {
	Barcode string  `json:"barcode,omitempty"`
	Events  []Event `json:"events,omitempty"`
}

type Event struct {
	Accepted  bool   `json:"accepted,omitempty"`
	Attribute string `json:"attribute,omitempty"`
	//	EventDate time.Time `json:"eventDate"`
	Location  string `json:"location"`
	Operation string `json:"operation"`
	Position  int    `json:"position"`
}

type Delivered struct {
	TrackerDelivered DeliverData `json:"trackerDelivered"`
}

type DeliverData struct {
	Barcode string `json:"barcode"`
}

type Message struct {
	Event   string `json:"event"`
	Barcode string `json:"barcode"`
	Status  string `json:"status"`
	//EventDate time.Time `json:"eventDate"`
}
