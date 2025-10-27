package main

import (
	"context"
	"fmt"
	"log"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	// "github.com/openai/openai-go/packages/param"
	// "github.com/openai/openai-go/responses"
)

func main() {
	client := openai.NewClient(option.WithBaseURL("http://127.0.0.1:9652/v1"), option.WithAPIKey("sk-dummy"))

	resp, err := client.Chat.Completions.New(context.Background(), openai.ChatCompletionNewParams{
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.UserMessage("Hello! Can you tell me a joke?"),
		},
		MaxTokens:   openai.Int(100),
		Temperature: openai.Float(4),
	})

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(resp.RawJSON())
	fmt.Println(resp.Choices[0].Message.Content)
}
