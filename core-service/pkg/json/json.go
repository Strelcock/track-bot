package json

import "encoding/json"

func Unmarshal[T any](data []byte) (*T, error) {
	var res T
	err := json.Unmarshal(data, &res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func Marshal[T any](src T) ([]byte, error) {
	data, err := json.Marshal(src)
	if err != nil {
		return nil, err
	}
	return data, nil
}
