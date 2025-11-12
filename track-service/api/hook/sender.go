package hook

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
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
		log.Println("try to get: ", err)
		return "", fmt.Errorf("carrier: %w", err)
	}
	if resp.StatusCode >= 400 {
		log.Println(resp.Status, resp.StatusCode)
		return "", fmt.Errorf("carrier: %d %s", resp.StatusCode, resp.Status)
	}
	err = json.NewDecoder(resp.Body).Decode(&carrier)
	if err != nil {
		return "", fmt.Errorf("carrier: %w", err)
	}
	if len(carrier) == 0 {
		return "", errors.New("carrier: not found")
	}

	return carrier[0].Code, nil
}

func (s *Sender) AddTracker(carrier, barcode string) error {
	url := fmt.Sprintf("https://moyaposylka.ru/api/v1/trackers/%s/%s", carrier, barcode)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return err
	}
	//log.Println(req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", s.apiKey)

	client := &http.Client{}
	//log.Println(client)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("sender: add tracker: ", err)
		return err
	}

	if resp.StatusCode >= 400 {
		var err HookAddError
		json.NewDecoder(resp.Body).Decode(&err)

		return err
	}

	return nil
}
