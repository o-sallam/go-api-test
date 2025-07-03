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

func BuildCardsHTML(articles []models.DummyArticle) string {
	var cardsBuilder strings.Builder
	for _, a := range articles {
		card := cardTemplate
		card = strings.ReplaceAll(card, "{{ALT}}", a.ALT)
		card = strings.ReplaceAll(card, "{{IMG}}", a.IMG)
		card = strings.ReplaceAll(card, "{{CATEGORY}}", a.CATEGORY)
		card = strings.ReplaceAll(card, "{{LINK}}", a.LINK)
		card = strings.ReplaceAll(card, "{{TITLE}}", a.TITLE)
		card = strings.ReplaceAll(card, "{{EXCERPT}}", a.EXCERPT)
		card = strings.ReplaceAll(card, "{{VIEWS}}", a.VIEWS)
		card = strings.ReplaceAll(card, "{{AUTHOR}}", a.AUTHOR)
		card = strings.ReplaceAll(card, "{{DATE}}", a.DATE)
		card = strings.ReplaceAll(card, "{{SLUG}}", a.Slug)
		cardsBuilder.WriteString(card)
	}
	return cardsBuilder.String()
}
