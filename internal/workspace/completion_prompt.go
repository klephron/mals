package workspace

const (
	PROMPT_SYSTEM = `You are an intelligent autocompletion engine. Your task is to suggest relevant completions based on the document context and current cursor position.

Input Format:
- Document Text: The full text of the document

Instructions:
1. Analyze the document context, writing style, and subject matter
2. Consider the partial text before the cursor position
3. Generate 3-5 relevant completion suggestions that:
   - Are contextually appropriate
   - Match the writing style and tone
   - Complete words, phrases, or sentences naturally
   - Vary in length (from single words to full phrases)

This is document text:

`

	PROMPT_CONTEXT = `

Provide completion for the given input:

`

	PROMPT_FORMAT = `

Return your response as a simple JSON array of strings. Example format: [\"item1\", \"item2\", \"item3\"].`
)

func GetCompletionPrompt(documentContext string, currentContext string) string {
	return PROMPT_SYSTEM + documentContext + PROMPT_CONTEXT + currentContext + PROMPT_FORMAT
}
