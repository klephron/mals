package core

type ModelSchema string

const (
	ModelSchemaJsonCompletionItem ModelSchema = "json/completionItem"
)

type ModelRole string

const (
	ModelRoleSystem    ModelRole = "system"
	ModelRoleUser      ModelRole = "user"
	ModelRoleAssistant ModelRole = "assistant"
)

type ModelMessage struct {
	Role    ModelRole
	Content string
}

type ModelParameters struct {
	Schema      ModelSchema
	Temperature *float64
	MaxTokens   *int64
}
