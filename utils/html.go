package utils

import "strings"

// MinifyHTML removes newlines and extra spaces between tags
func MinifyHTML(input string) string {
	return strings.Join(strings.Fields(input), " ")
}
