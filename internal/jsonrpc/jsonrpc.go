package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"mals-engine/pkg/message"
	"strconv"
)

func getSeparator() []byte {
	return []byte{'\r', '\n', '\r', '\n'}
}

func cutRequestMessage(data []byte) (header []byte, content []byte, err error) {
	header, content, found := bytes.Cut(data, getSeparator())
	if !found {
		err = errors.New("separator not found")
	}
	return
}

func getContentLength(header []byte) (length int, err error) {
	lengthBytes := header[len("Content-Length: "):]
	length, err = strconv.Atoi(string(lengthBytes))
	return
}

func DecodeRequest(data []byte) (message.Request, []byte, error) {
	header, content, err := cutRequestMessage(data)
	if err != nil {
		return message.Request{}, nil, err
	}

	length, err := getContentLength(header)
	if err != nil {
		return message.Request{}, nil, err
	}

	var requestMessage message.Request
	if err := json.Unmarshal(content[:length], &requestMessage); err != nil {
		return message.Request{}, nil, err
	}

	return requestMessage, content[:length], nil
}

func EncodeResponse(msg any) ([]byte, error) {
	content, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	return fmt.Appendf(nil, "Content-Length: %d\r\n\r\n%s", len(content), content), nil
}

func ScannerSplit(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, err := cutRequestMessage(data)
	if err != nil {
		return 0, nil, nil
	}

	length, err := getContentLength(header)
	if err != nil {
		return 0, nil, err
	}

	if len(content) < length {
		return 0, nil, nil
	}

	totalLength := len(header) + len(getSeparator()) + length

	return totalLength, data[:totalLength], nil
}
