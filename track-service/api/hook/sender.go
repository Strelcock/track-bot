package hook

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Sender struct {
}

func NewSender() *Sender {
	return &Sender{}
}

func (s *Sender) Carrier(num string) (string, error) {
	var carrier CarrierResponse

	url := fmt.Sprintf("https://moyaposylka.ru/api/v1/carriers/%s", num)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("carrier: %w", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&carrier)
	if err != nil {
		return "", fmt.Errorf("carrier: %w", err)
	}
	return carrier.Code, nil
}
