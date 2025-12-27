package lsp

import "encoding/json"

func rawDecode[T any](s *ClientLsp, bytes json.RawMessage) (res T, err error) {
	err = json.Unmarshal(bytes, &res)
	if err != nil {
		s.plane.Log().Warnf("%v", err)
	}
	return
}

func rawEncode[T any](s *ClientLsp, data *T) (json.RawMessage, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		s.plane.Log().Errorf("%v", err)
		return nil, err
	}
	return raw, nil
}
