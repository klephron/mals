package jsonrpc

import (
	"bytes"
	"encoding/json"
	"errors"
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

func DecodeRequestMessage(data []byte) (message.RequestMessage, []byte, error) {
	header, content, err := cutRequestMessage(data)
	if err != nil {
		return message.RequestMessage{}, nil, err
	}

	length, err := getContentLength(header)
	if err != nil {
		return message.RequestMessage{}, nil, err
	}

	var requestMessage message.RequestMessage
	if err := json.Unmarshal(content[:length], &requestMessage); err != nil {
		return message.RequestMessage{}, nil, err
	}

	return requestMessage, content[:length], nil
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
