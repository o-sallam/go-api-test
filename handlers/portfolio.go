package handlers

import (
	"go-api-test/models"
	"net/http"
	"os"
	"strings"
)

var (
	portfolioHTML string
	cardTemplate  string
)

// SetPortfolioHTML sets the HTML template in memory
func SetPortfolioHTML(html string) {
	portfolioHTML = html
	// Load the card template at startup
	cardBytes, err := os.ReadFile("wwwroot/components/card.html")
	if err != nil {
		panic("Failed to load card template: " + err.Error())
	}
	cardTemplate = string(cardBytes)
}

// PortfolioHandler serves the HTML with dynamic cards
func PortfolioHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	// Generate dummy data
	articles := []models.DummyArticle{
		{
			ALT:      "كل ما تريد معرفته عن القهوة - صورة توضيحية",
			IMG:      "/img/blog.webp",
			CATEGORY: "تقنية",
			LINK:     "posts/first-blog-post/",
			TITLE:    "كل ما تريد معرفته عن القهوة",
			EXCERPT:  "تعرّف على أحدث تطبيقات الهواتف الذكية وكيفية استخدامها لتحسين الإنتاجية والترفيه.",
			VIEWS:    "1.5k",
			AUTHOR:   "محمد علي",
			DATE:     "منذ 4 أيام",
		},
		{
			ALT:      "Second Blog Post - صورة توضيحية",
			IMG:      "/img/blog.webp",
			CATEGORY: "صحة",
			LINK:     "posts/second-blog-post/",
			TITLE:    "Second Blog Post",
			EXCERPT:  "اكتشف أهمية النوم الصحي وكيف يؤثر على صحتك العامة ونشاطك اليومي.",
			VIEWS:    "1.2k",
			AUTHOR:   "محمد علي",
			DATE:     "منذ 3 أيام",
		},
	}

	// Build cards HTML
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
		cardsBuilder.WriteString(card)
	}

	// Inject cards into the main HTML
	out := strings.Replace(portfolioHTML, "{{CARDS}}", cardsBuilder.String(), 1)
	w.Write([]byte(out))
}
