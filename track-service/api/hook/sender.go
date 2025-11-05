package hook

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

var ()

type Sender struct {
	apiKey string
}

func NewSender(apiKey string) *Sender {
	return &Sender{apiKey: apiKey}
}

func (s *Sender) Carrier(barcode string) (string, error) {
	var carrier []CarrierResponse

	url := fmt.Sprintf("https://moyaposylka.ru/api/v1/carriers/%s", barcode)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("carrier: %w", err)
	}
	err = json.NewDecoder(resp.Body).Decode(&carrier)
	if err != nil {
		return "", fmt.Errorf("carrier: %w", err)
	}
	return carrier[0].Code, nil
}

func (s *Sender) AddTracker(carrier, barcode string) error {
	url := fmt.Sprintf("https://moyaposylka.ru/api/v1/trackers/%s/%s", carrier, barcode)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.apiKey)

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("%d, %s", resp.StatusCode, resp.Status)
	}

	return nil
}
