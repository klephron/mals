package jsonrpc

func ScannerSplit(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, err := cutParts(data)
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
