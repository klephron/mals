package model

import (
	"context"
	"errors"
	"log"
	"mals-engine/pkg/config"
)

type ModelRequest struct {
	Text string
}

type ModelResponse struct {
	Text  string
	Error error
}

type ModelService interface {
	Serve(ctx context.Context)
	onRequest(ctx context.Context, request ModelRequest) ModelResponse
	NewModelRequest(request ModelRequest) ModelResponse
}

type Model struct {
	ModelService
	logger    *log.Logger
	requests  chan ModelRequest
	responses chan ModelResponse
	Id        string
	Spec      string
	BaseUrl   string
	Settings  config.ModelSettings
}

func NewModelResponse(text string) ModelResponse {
	return ModelResponse{Text: text, Error: nil}
}

func NewModelError(error error) ModelResponse {
	return ModelResponse{Text: "", Error: error}
}

// current client must wait until response, otherwise model will hang
func (m *Model) NewModelRequest(request ModelRequest) ModelResponse {
	m.requests <- request
	return <-m.responses
}

func (m *Model) Serve(ctx context.Context) {
	m.LogInfoPrintf("serving started")
	for {
		select {
		case <-ctx.Done():
			close(m.requests)
			close(m.responses)
			m.LogInfoPrintf("serving stopped")
			return
		case request := <-m.requests:
			m.responses <- m.ModelService.onRequest(ctx, request)
		}
	}
}

func (m *Model) onRequest(ctx context.Context, request ModelRequest) ModelResponse {
	return NewModelError(errors.New("generic model is unable to generate responses"))
}
