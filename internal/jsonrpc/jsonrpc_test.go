package jsonrpc

import (
	"bytes"
	"testing"
)

func TestScannerSplit(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		expectedAdvance int
		expectedToken   []byte
		expectedError   bool
	}{
		{
			name:            "complete valid message",
			input:           []byte("Content-Length: 17\r\n\r\n{\"jsonrpc\":\"2.0\"}"),
			expectedAdvance: 39,
			expectedToken:   []byte("Content-Length: 17\r\n\r\n{\"jsonrpc\":\"2.0\"}"),
			expectedError:   false,
		},
		{
			name:            "incomplete message no separator",
			input:           []byte("Content-Length: 17"),
			expectedAdvance: 0,
			expectedToken:   nil,
			expectedError:   false,
		},
		{
			name:            "incomplete message partial content",
			input:           []byte("Content-Length: 17\r\n\r\n{\"jsonrpc\""),
			expectedAdvance: 0,
			expectedToken:   nil,
			expectedError:   false,
		},
		{
			name:            "zero length content",
			input:           []byte("Content-Length: 0\r\n\r\n"),
			expectedAdvance: 21,
			expectedToken:   []byte("Content-Length: 0\r\n\r\n"),
			expectedError:   false,
		},
		{
			name:            "invalid content length",
			input:           []byte("Content-Length: abc\r\n\r\n{\"test\": true}"),
			expectedAdvance: 0,
			expectedToken:   nil,
			expectedError:   true,
		},
		{
			name:            "exact content length match",
			input:           []byte("Content-Length: 5\r\n\r\nhello"),
			expectedAdvance: 26,
			expectedToken:   []byte("Content-Length: 5\r\n\r\nhello"),
			expectedError:   false,
		},
		{
			name:            "content longer than specified",
			input:           []byte("Content-Length: 5\r\n\r\nhello world"),
			expectedAdvance: 26,
			expectedToken:   []byte("Content-Length: 5\r\n\r\nhello"),
			expectedError:   false,
		},
		{
			name:            "multiple messages in buffer",
			input:           []byte("Content-Length: 5\r\n\r\nhelloContent-Length: 5\r\n\r\nworld"),
			expectedAdvance: 26,
			expectedToken:   []byte("Content-Length: 5\r\n\r\nhello"),
			expectedError:   false,
		},
		{
			name:            "empty input",
			input:           []byte(""),
			expectedAdvance: 0,
			expectedToken:   nil,
			expectedError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			advance, token, err := ScannerSplit(tt.input, false)

			if tt.expectedError {
				if err == nil {
					t.Errorf("ScannerSplit() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ScannerSplit() unexpected error: %v", err)
				return
			}

			if advance != tt.expectedAdvance {
				t.Errorf("ScannerSplit() advance = %d, want %d", advance, tt.expectedAdvance)
			}

			if !bytes.Equal(token, tt.expectedToken) {
				t.Errorf("ScannerSplit() token = %q, want %q", token, tt.expectedToken)
			}
		})
	}
}

func TestDecodeMessageRequest(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expected      *Request
		expectedError bool
	}{
		{
			name:  "valid request",
			input: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":1}"),
			expected: &Request{
				Id:     1,
				Method: "test",
				Params: nil,
			},
			expectedError: false,
		},
		{
			name:  "request with different method",
			input: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"ping\",\"id\":2}"),
			expected: &Request{
				Id:     2,
				Method: "ping",
				Params: nil,
			},
			expectedError: false,
		},
		{
			name:          "no separator",
			input:         []byte("Content-Length: 38"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid content length",
			input:         []byte("Content-Length: abc\r\n\r\n{\"jsonrpc\":\"2.0\"}"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid JSON",
			input:         []byte("Content-Length: 15\r\n\r\n{\"invalid json\""),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "negative ID",
			input: []byte("Content-Length: 41\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":-1}"),
			expected: &Request{
				Method: "test",
				Id:     -1,
			},
			expectedError: false,
		},
		{
			name:  "zero ID",
			input: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":0}"),
			expected: &Request{
				Method: "test",
				Id:     0,
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := DecodeMessage(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("DecodeMessage() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeMessage() unexpected error: %v", err)
				return
			}

			switch req := msg.(type) {
			case *Request:
				if req.Method != tt.expected.Method {
					t.Errorf("DecodeMessage() method = %q, want %q", req.Method, tt.expected.Method)
				}

				if req.Id != tt.expected.Id {
					t.Errorf("DecodeRequest() id = %d, want %d", req.Id, tt.expected.Id)
				}
			default:
				t.Errorf("DecodeMessage() expected %T, got %T", &Request{}, msg)
			}

		})
	}
}

func TestDecodeMessageNotification(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expected      *Notification
		expectedError bool
	}{
		{
			name:  "valid notification",
			input: []byte("Content-Length: 35\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"notify\"}"),
			expected: &Notification{
				Method: "notify",
			},
			expectedError: false,
		},
		{
			name:  "notification with different method",
			input: []byte("Content-Length: 33\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"ping\"}"),
			expected: &Notification{
				Method: "ping",
			},
			expectedError: false,
		},
		{
			name:          "no separator",
			input:         []byte("Content-Length: 33"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid JSON",
			input:         []byte("Content-Length: 15\r\n\r\n{\"invalid json\""),
			expected:      nil,
			expectedError: true,
		},
		{
			name:          "invalid content length",
			input:         []byte("Content-Length: abc\r\n\r\n{\"jsonrpc\":\"2.0\"}"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "long method name",
			input: []byte("Content-Length: 49\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"textDocument/didOpen\"}"),
			expected: &Notification{
				Method: "textDocument/didOpen",
			},
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := DecodeMessage(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("DecodeMessage() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeMessage() unexpected error: %v", err)
				return
			}

			switch req := msg.(type) {
			case *Notification:
				if req.Method != tt.expected.Method {
					t.Errorf("DecodeMessage() method = %q, want %q", req.Method, tt.expected.Method)
				}
			default:
				t.Errorf("DecodeMessage() expected %T, got %T", &Notification{}, msg)
			}
		})
	}
}

func TestDecodeMessageResponse(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expected      *Response
		expectedError bool
	}{
		{
			name:  "valid response with result",
			input: []byte("Content-Length: 38\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":1,\"result\":true}"),
			expected: &Response{
				Id: 1,
			},
			expectedError: false,
		},
		{
			name:  "valid response with error",
			input: []byte("Content-Length: 48\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":2,\"error\":{\"code\":-32600}}"),
			expected: &Response{
				Id: 2,
			},
			expectedError: false,
		},
		{
			name:  "response with null result",
			input: []byte("Content-Length: 38\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":3,\"result\":null}"),
			expected: &Response{
				Id: 3,
			},
			expectedError: false,
		},
		{
			name:  "response with large id",
			input: []byte("Content-Length: 47\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":2147483647,\"result\":\"ok\"}"),
			expected: &Response{
				Id: 2147483647,
			},
			expectedError: false,
		},
		{
			name:          "response missing id",
			input:         []byte("Content-Length: 31\r\n\r\n{\"jsonrpc\":\"2.0\",\"result\":true}"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "response with both result and error",
			input: []byte("Content-Length: 58\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":4,\"result\":true,\"error\":{\"code\":-1}}"),
			expected: &Response{
				Id: 4,
			},
			expectedError: false,
		},
		{
			name:  "response with complex result",
			input: []byte("Content-Length: 49\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":5,\"result\":{\"key\":\"value\"}}"),
			expected: &Response{
				Id: 5,
			},
			expectedError: false,
		},
		{
			name:          "invalid json in response",
			input:         []byte("Content-Length: 22\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":"),
			expected:      nil,
			expectedError: true,
		},
		{
			name:  "response with zero id",
			input: []byte("Content-Length: 38\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":0,\"result\":true}"),
			expected: &Response{
				Id: 0,
			},
			expectedError: false,
		},
		{
			name:          "no separator",
			input:         []byte("Content-Length: 44"),
			expected:      nil,
			expectedError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg, err := DecodeMessage(tt.input)
			if tt.expectedError {
				if err == nil {
					t.Errorf("DecodeMessage() expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("DecodeMessage() unexpected error: %v", err)
				return
			}
			switch resp := msg.(type) {
			case *Response:
				if resp.Id != tt.expected.Id {
					t.Errorf("DecodeMessage() id = %d, want %d", resp.Id, tt.expected.Id)
				}
			default:
				t.Errorf("DecodeMessage() expected %T, got %T", &Response{}, msg)
			}
		})
	}
}

func TestEncodeMessage(t *testing.T) {
	tests := []struct {
		name           string
		input          Message
		expectedOutput []byte
		expectedError  bool
	}{
		{
			name: "encode request",
			input: &Request{
				Method: "test",
				Id:     1,
			},
			expectedOutput: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":1,\"method\":\"test\"}"),
			expectedError:  false,
		},
		{
			name: "encode notification",
			input: &Notification{
				Method: "notify",
			},
			expectedOutput: []byte("Content-Length: 35\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"notify\"}"),
			expectedError:  false,
		},
		{
			name: "encode response",
			input: &Response{
				Id: 42,
			},
			expectedOutput: []byte("Content-Length: 25\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":42}"),
			expectedError:  false,
		},
		{
			name: "encode request with negative ID",
			input: &Request{
				Method: "test",
				Id:     -1,
			},
			expectedOutput: []byte("Content-Length: 41\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":-1,\"method\":\"test\"}"),
			expectedError:  false,
		},
		{
			name: "encode LSP method",
			input: &Request{
				Id:     123,
				Method: "textDocument/hover",
			},
			expectedOutput: []byte("Content-Length: 56\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":123,\"method\":\"textDocument/hover\"}"),
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := EncodeMessage(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Encode() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Encode() unexpected error: %v", err)
				return
			}

			if !bytes.Equal(result, tt.expectedOutput) {
				t.Errorf("Encode() = %q, want %q", result, tt.expectedOutput)
			}
		})
	}
}
