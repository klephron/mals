package test

import (
	"bytes"
	. "mals-engine/internal/jsonrpc"
	. "mals-engine/pkg/lsp_message"
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

func TestDecodeRequest(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expectedReq   Request
		expectedRaw   []byte
		expectedError bool
	}{
		{
			name:  "valid request",
			input: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":1}"),
			expectedReq: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "test",
				},
				Id: 1,
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":1}"),
			expectedError: false,
		},
		{
			name:  "request with different method",
			input: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"ping\",\"id\":2}"),
			expectedReq: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "ping",
				},
				Id: 2,
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"ping\",\"id\":2}"),
			expectedError: false,
		},
		{
			name:          "no separator",
			input:         []byte("Content-Length: 38"),
			expectedReq:   Request{},
			expectedRaw:   nil,
			expectedError: true,
		},
		{
			name:          "invalid content length",
			input:         []byte("Content-Length: abc\r\n\r\n{\"jsonrpc\":\"2.0\"}"),
			expectedReq:   Request{},
			expectedRaw:   nil,
			expectedError: true,
		},
		{
			name:          "invalid JSON",
			input:         []byte("Content-Length: 15\r\n\r\n{\"invalid json\""),
			expectedReq:   Request{},
			expectedRaw:   nil,
			expectedError: true,
		},
		{
			name:  "negative ID",
			input: []byte("Content-Length: 41\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":-1}"),
			expectedReq: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "test",
				},
				Id: -1,
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":-1}"),
			expectedError: false,
		},
		{
			name:  "zero ID",
			input: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":0}"),
			expectedReq: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "test",
				},
				Id: 0,
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":0}"),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, raw, err := DecodeRequest(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("DecodeRequest() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeRequest() unexpected error: %v", err)
				return
			}

			if req.JsonRpc != tt.expectedReq.JsonRpc {
				t.Errorf("DecodeRequest() jsonrpc = %q, want %q", req.JsonRpc, tt.expectedReq.JsonRpc)
			}

			if req.Method != tt.expectedReq.Method {
				t.Errorf("DecodeRequest() method = %q, want %q", req.Method, tt.expectedReq.Method)
			}

			if req.Id != tt.expectedReq.Id {
				t.Errorf("DecodeRequest() id = %d, want %d", req.Id, tt.expectedReq.Id)
			}

			if !bytes.Equal(raw, tt.expectedRaw) {
				t.Errorf("DecodeRequest() raw = %q, want %q", raw, tt.expectedRaw)
			}
		})
	}
}

func TestDecodeNotification(t *testing.T) {
	tests := []struct {
		name          string
		input         []byte
		expectedNotif Notification
		expectedRaw   []byte
		expectedError bool
	}{
		{
			name:  "valid notification",
			input: []byte("Content-Length: 35\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"notify\"}"),
			expectedNotif: Notification{
				Message: Message{JsonRpc: "2.0"},
				Method:  "notify",
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"notify\"}"),
			expectedError: false,
		},
		{
			name:  "notification with different method",
			input: []byte("Content-Length: 33\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"ping\"}"),
			expectedNotif: Notification{
				Message: Message{JsonRpc: "2.0"},
				Method:  "ping",
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"ping\"}"),
			expectedError: false,
		},
		{
			name:          "no separator",
			input:         []byte("Content-Length: 33"),
			expectedNotif: Notification{},
			expectedRaw:   nil,
			expectedError: true,
		},
		{
			name:          "invalid JSON",
			input:         []byte("Content-Length: 15\r\n\r\n{\"invalid json\""),
			expectedNotif: Notification{},
			expectedRaw:   nil,
			expectedError: true,
		},
		{
			name:          "invalid content length",
			input:         []byte("Content-Length: abc\r\n\r\n{\"jsonrpc\":\"2.0\"}"),
			expectedNotif: Notification{},
			expectedRaw:   nil,
			expectedError: true,
		},
		{
			name:  "long method name",
			input: []byte("Content-Length: 49\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"textDocument/didOpen\"}"),
			expectedNotif: Notification{
				Message: Message{JsonRpc: "2.0"},
				Method:  "textDocument/didOpen",
			},
			expectedRaw:   []byte("{\"jsonrpc\":\"2.0\",\"method\":\"textDocument/didOpen\"}"),
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			notif, raw, err := DecodeNotification(tt.input)

			if tt.expectedError {
				if err == nil {
					t.Errorf("DecodeNotification() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("DecodeNotification() unexpected error: %v", err)
				return
			}

			if notif.JsonRpc != tt.expectedNotif.JsonRpc {
				t.Errorf("DecodeNotification() jsonrpc = %q, want %q", notif.JsonRpc, tt.expectedNotif.JsonRpc)
			}

			if notif.Method != tt.expectedNotif.Method {
				t.Errorf("DecodeNotification() method = %q, want %q", notif.Method, tt.expectedNotif.Method)
			}

			if !bytes.Equal(raw, tt.expectedRaw) {
				t.Errorf("DecodeNotification() raw = %q, want %q", raw, tt.expectedRaw)
			}
		})
	}
}

func TestEncode(t *testing.T) {
	tests := []struct {
		name           string
		input          any
		expectedOutput []byte
		expectedError  bool
	}{
		{
			name: "encode request",
			input: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "test",
				},
				Id: 1,
			},
			expectedOutput: []byte("Content-Length: 40\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":1}"),
			expectedError:  false,
		},
		{
			name: "encode notification",
			input: Notification{
				Message: Message{JsonRpc: "2.0"},
				Method:  "notify",
			},
			expectedOutput: []byte("Content-Length: 35\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"notify\"}"),
			expectedError:  false,
		},
		{
			name: "encode response",
			input: Response{
				Message: Message{JsonRpc: "2.0"},
				Id:      42,
			},
			expectedOutput: []byte("Content-Length: 25\r\n\r\n{\"jsonrpc\":\"2.0\",\"id\":42}"),
			expectedError:  false,
		},
		{
			name: "encode request with negative ID",
			input: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "test",
				},
				Id: -1,
			},
			expectedOutput: []byte("Content-Length: 41\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"test\",\"id\":-1}"),
			expectedError:  false,
		},
		{
			name: "encode LSP method",
			input: Request{
				Notification: Notification{
					Message: Message{JsonRpc: "2.0"},
					Method:  "textDocument/hover",
				},
				Id: 123,
			},
			expectedOutput: []byte("Content-Length: 56\r\n\r\n{\"jsonrpc\":\"2.0\",\"method\":\"textDocument/hover\",\"id\":123}"),
			expectedError:  false,
		},
		{
			name:          "unencodable input",
			input:         make(chan int), // channels can't be JSON marshaled
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := Encode(tt.input)

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
