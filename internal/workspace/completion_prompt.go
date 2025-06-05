package workspace

import (
	"fmt"
)

const (
	PROMPT_TEMPLATE = `
You are an intelligent autocompletion engine. Your task is to suggest relevant completions based on the document context and current cursor position.

Input Format:
- Document Text: The full text of the document
- Current Context: Text before current position

Instructions:
1. Analyze the document context, writing style, and subject matter
2. Consider the partial text before the cursor position
3. No more than 5 relevant completion suggestions that:
   - Are contextually appropriate
   - Match the writing style and tone
   - Complete words, phrases, or sentences naturally
   - Vary in length (from single words to full phrases)

Document text:

%s

Return your response as a simple JSON array of strings without duplicates. Example format: [\"item1\", \"item2\", \"item3\"].

Complete this text:

%s`

)

func GetCompletionPrompt(documentContext string, currentContext string) string {
	prompt := fmt.Sprintf(PROMPT_TEMPLATE, documentContext, currentContext)
	fmt.Println(prompt)
	return prompt
}
