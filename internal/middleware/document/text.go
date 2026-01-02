package document

import "strings"

func textToLines(text *string) []string {
	if text != nil {
		return strings.Split(*text, "\n")
	}
	return []string{}
}
