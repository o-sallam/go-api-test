package utils

import (
	"github.com/tdewolff/minify/v2"
	cssmin "github.com/tdewolff/minify/v2/css"
	htmlmin "github.com/tdewolff/minify/v2/html"
)

// MinifyHTML uses tdewolff/minify for HTML minification
func MinifyHTML(input string) string {
	m := minify.New()
	m.Add("text/html", &htmlmin.Minifier{})
	minified, err := m.String("text/html", input)
	if err != nil {
		return input // fallback to original if error
	}
	return minified
}

// MinifyCSS uses tdewolff/minify for CSS minification
func MinifyCSS(input string) string {
	m := minify.New()
	m.Add("text/css", &cssmin.Minifier{})
	minified, err := m.String("text/css", input)
	if err != nil {
		return input // fallback to original if error
	}
	return minified
}
