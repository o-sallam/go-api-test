package utils

import "strings"

// MinifyHTML removes newlines and extra spaces between tags
func MinifyHTML(input string) string {
	return strings.Join(strings.Fields(input), " ")
}

// MinifyCSS removes newlines and extra spaces from CSS
func MinifyCSS(input string) string {
	return strings.Join(strings.Fields(input), " ")
}
