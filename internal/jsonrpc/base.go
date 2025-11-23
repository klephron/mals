package jsonrpc

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
)

func getSeparator() []byte {
	return []byte{'\r', '\n', '\r', '\n'}
}

func cutParts(data []byte) (header []byte, content []byte, err error) {
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

func decode(data []byte) (body []byte, err error) {
	header, content, err := cutParts(data)
	if err != nil {
		return nil, err
	}

	length, err := getContentLength(header)
	if err != nil {
		return nil, err
	}

	if length != len(content) {
		return nil, fmt.Errorf("decode content length: expected %v, got %v", len(content), length)
	}

	return content[:length], nil
}

func encode(body []byte) ([]byte, error) {
	return fmt.Appendf(nil, "Content-Length: %d\r\n\r\n%s", len(body), body), nil
}
