package utils

import (
	"go-api-test/models"
	"os"
	"strings"
)

var cardTemplate string

func LoadCardTemplate(path string) error {
	cardBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	cardTemplate = string(cardBytes)
	return nil
}

// ReplacePlaceholders replaces all {{KEY}} in template with values from fields map
func ReplacePlaceholders(template string, fields map[string]string) string {
	for key, value := range fields {
		template = strings.ReplaceAll(template, "{{"+key+"}}", value)
	}
	return template
}

func BuildCardsHTML(articles []models.PostCardResponse) string {
	var cardsBuilder strings.Builder
	for _, a := range articles {
		fields := map[string]string{
			"ALT":      a.ALT,
			"IMG":      a.IMG,
			"CATEGORY": a.CATEGORY,
			"LINK":     a.LINK,
			"TITLE":    a.TITLE,
			"EXCERPT":  a.EXCERPT,
			"VIEWS":    a.VIEWS,
			"AUTHOR":   a.AUTHOR,
			"DATE":     a.DATE,
			"SLUG":     a.Slug,
			// يمكنك إضافة أي حقول ديناميكية أخرى هنا
		}
		card := ReplacePlaceholders(cardTemplate, fields)
		cardsBuilder.WriteString(card)
	}
	return cardsBuilder.String()
}
