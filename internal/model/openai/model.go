package openai

import (
	"context"
	"mals/internal/model"
	"mals/internal/plane"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/param"
)

type ModelOpenAISpec struct {
	Url         string
	MaxTokens   int
	Temperature float32
}

type ModelOpenAI struct {
	model.Model
	name string

	client      openai.Client
	maxTokens   param.Opt[int64]
	temperature param.Opt[float64]

	plane plane.Plane
}

func New(name string, spec ModelOpenAISpec, plane plane.Plane) (*ModelOpenAI, error) {
	client := openai.NewClient(option.WithBaseURL(spec.Url), option.WithAPIKey("sk-dummy"))
	maxTokens := openai.Int(int64(spec.MaxTokens))
	temperature := openai.Float(float64(spec.Temperature))

	return &ModelOpenAI{
		name:        name,
		client:      client,
		maxTokens:   maxTokens,
		temperature: temperature,
		plane:       plane,
	}, nil
}

func (s *ModelOpenAI) Name() string {
	return s.name
}

func (s *ModelOpenAI) Kind() string {
	return Kind()
}

func (s *ModelOpenAI) Run(ctx context.Context) error {
	s.plane.Infof("%T %v: started", s, s.Name())

	<-ctx.Done()

	s.plane.Infof("%T %v: done", s, s.Name())

	return nil
}

func (s *ModelOpenAI) Execute(ctx context.Context, task *model.Task) (string, error) {

	s.plane.Infof("%T %v task %v: received", s, s.Name(), task)

	schema := openai.ResponseFormatJSONSchemaJSONSchemaParam{
		Name:   task.SchemaName,
		Schema: task.Schema,
		Strict: openai.Bool(task.SchemaStrict),
	}

	resp, err := s.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage(task.Text),
		},
		MaxTokens:   s.maxTokens,
		Temperature: s.temperature,
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: schema,
			},
		},
	})

	s.plane.Infof("%T %v task %v: processed", s, s.Name(), task)

	if err != nil {
		s.plane.Warnf("%T %v task %v: ", s, s.Name(), task.Id, err)
		return "", err
	}

	text := resp.Choices[0].Message.Content

	return text, nil
}
