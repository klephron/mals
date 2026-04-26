package openai

import (
	"context"
	"mals/internal/model"
	"mals/internal/plane"
	"mals/pkg/config"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

type ModelOpenai struct {
	name string

	client      openai.Client
	maxTokens   param.Opt[int64]
	temperature param.Opt[float64]

	plane plane.Plane
}

func New(name string, api *config.ModelApiOpenai, plane plane.Plane) *ModelOpenai {
	client := openai.NewClient(option.WithBaseURL(api.Url), option.WithAPIKey("sk-dummy"))

	var maxTokens param.Opt[int64]
	if api.MaxTokens != nil {
		maxTokens = openai.Int(int64(*api.MaxTokens))
	}

	var temperature param.Opt[float64]
	if api.Temperature != nil {
		temperature = openai.Float(float64(*api.Temperature))
	}

	return &ModelOpenai{
		name:        name,
		client:      client,
		maxTokens:   maxTokens,
		temperature: temperature,
		plane:       plane,
	}
}

func (s *ModelOpenai) Name() string {
	return s.name
}

func (s *ModelOpenai) Kind() string {
	var settings config.ModelApiOpenai
	return settings.ModelApiKind()
}

func (s *ModelOpenai) Run(ctx context.Context) error {
	s.plane.Infof("%T %v: started", s, s.Name())

	<-ctx.Done()

	s.plane.Infof("%T %v: done", s, s.Name())

	return nil
}

func (s *ModelOpenai) Execute(ctx context.Context, task *model.Task) (string, error) {

	s.plane.Infof("%T %v task %v: received", s, s.Name(), task)

	var responseFormat openai.ChatCompletionNewParamsResponseFormatUnion

	switch task.Schema {
	case config.StepModelSchemaJsonCompletionItems:
		completionItemsSchema := map[string]any{
			"type": "array",
			"items": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"label": map[string]any{
						"type": "string",
					},
				},
				"required":             []string{"label"},
				"additionalProperties": false,
			},
		}
		responseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   string(config.StepModelSchemaJsonCompletionItems),
					Strict: openai.Bool(true),
					Schema: completionItemsSchema,
				},
			},
		}
	default:
		responseFormat = openai.ChatCompletionNewParamsResponseFormatUnion{}
	}

	resp, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(task.Text),
		},
		MaxCompletionTokens: s.maxTokens,
		Temperature:         s.temperature,
		N:                   openai.Int(1),
		ResponseFormat:      responseFormat,
	})

	s.plane.Infof("%T %v task %v: processed", s, s.Name(), task)

	if err != nil {
		s.plane.Warnf("%T %v task %v: ", s, s.Name(), task.Id, err)
		return "", err
	}

	text := resp.Choices[0].Message.Content

	return text, nil
}
