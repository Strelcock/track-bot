package hook

type CarrierResponse struct {
	Code string `json:"code"`
}

type NewEvent struct {
	Barcode string   `json:"barcode"`
	Events  []string `json:"events"`
}

type TrackDelivered struct {
	Barcode string `json:"barcode"`
}
