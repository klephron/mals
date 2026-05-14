package core

type ModelSchema string

const (
	ModelSchemaJsonCompletionItems ModelSchema = "json/completionItems"
)

type ModelRole string

const (
	ModelRoleSystem    ModelRole = "system"
	ModelRoleUser      ModelRole = "user"
	ModelRoleAssistant ModelRole = "assistant"
)

type ModelMessage struct {
	Role    ModelRole `json:"role"`
	Content string    `json:"content"`
}
