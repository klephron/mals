package workspace

// import (
// 	"fmt"
// )

// const (
// 	PROMPT_TEMPLATE = `
// <|begin_of_text|><|start_header_id|>system<|end_header_id|>

// You are an intelligent autocompletion engine. Your task is to suggest relevant text completions based on document context.

// Rules:
// - Analyze the document context and writing style
// - Generate 3-5 contextually appropriate completions
// - Match the existing tone and style
// - Provide completions of varying lengths (words to phrases)
// - Return ONLY a valid JSON array of strings
// - No explanations, no additional text

// <|eot_id|><|start_header_id|>user<|end_header_id|>

// Document text: The weather today is quite nice. I think I'll go for a walk in the park and maybe stop by the
// Current context: "I think I'll go for a walk in the park and maybe stop by the"

// <|eot_id|><|start_header_id|>assistant<|end_header_id|>

// ["coffee shop", "library", "store to buy groceries", "lake to feed the ducks"]

// <|eot_id|><|start_header_id|>user<|end_header_id|>

// Document text: She opened her laptop and started typing her report. The first chapter discussed the importance of renewable energy sources. Solar panels and wind turbines are becoming more
// Current context: "Solar panels and wind turbines are becoming more"

// <|eot_id|><|start_header_id|>assistant<|end_header_id|>

// ["efficient", "affordable each year", "popular worldwide", "cost-effective solutions", "widely adopted"]

// <|eot_id|><|start_header_id|>user<|end_header_id|>

// Document text:
// %s

// Current context:
// %s

// <|eot_id|><|start_header_id|>assistant<|end_header_id|>
// `
// )

// func GetCompletionPrompt(documentContext string, currentContext string) string {
// 	prompt := fmt.Sprintf(PROMPT_TEMPLATE, documentContext, currentContext)
// 	// fmt.Println(prompt)
// 	return prompt
// }
