package hook

type CarrierResponse struct {
	Code string `json:"code"`
}

type Event struct {
	NewEvents EventData `json:"newEvents"`
}

type EventData struct {
	Barcode string   `json:"barcode,omitempty"`
	Events  []string `json:"events,omitempty"`
}

type Delivered struct {
	TrackerDelivered DeliverData `json:"trackerDelivered"`
}

type DeliverData struct {
	Barcode string `json:"barcode"`
}
