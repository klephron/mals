package model

import (
	"context"
	"log"
	"mals/pkg/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

type ModelOpenAI struct {
	Model
	client        openai.Client
	clientCreated bool
}

func NewModelOpenAI(logger *log.Logger, id string, spec string, baseUrl string, settings config.ModelSettings) *ModelOpenAI {
	m := &ModelOpenAI{
		Model:         NewModel(logger, id, spec, baseUrl, settings),
		client:        openai.Client{},
		clientCreated: false,
	}
	m.ModelService = m
	return m
}

func (m *ModelOpenAI) onRequest(ctx context.Context, request ModelRequest) ModelResponse {
	if !m.clientCreated {
		m.client = openai.NewClient(option.WithBaseURL(m.BaseUrl), option.WithAPIKey("sk-dummy"))
		m.clientCreated = true
	}

	var maxTokens param.Opt[int64]
	if m.Settings.MaxTokens != nil {
		maxTokens = openai.Int(*m.Settings.MaxTokens)
	}

	var temperature param.Opt[float64]
	if m.Settings.Temperature != nil {
		temperature = openai.Float(*m.Settings.Temperature)
	}

	schemaParam := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   request.SchemaName,
		Schema: request.Schema,
		Strict: openai.Bool(true),
	}

	resp, err := m.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(request.Text),
		},
		MaxTokens:   maxTokens,
		Temperature: temperature,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schemaParam,
			},
		},
	})

	if err != nil {
		return NewModelError(err)
	}
	return NewModelResponse(resp.Choices[0].Message.Content)
}
