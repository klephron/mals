package util

import "encoding/json"

func Ptr[T any](s T) *T {
	return &s
}

func JsonUnmarshal[T any](bytes json.RawMessage) (res T, err error) {
	err = json.Unmarshal(bytes, &res)
	return
}

func JsonMarshal[T any](data *T) (json.RawMessage, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
