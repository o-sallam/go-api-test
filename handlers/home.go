package handlers

import (
	"go-api-test/models"
	"go-api-test/utils"
	"net/http"
	"strings"
)

var portfolioHTML string

// SetPortfolioHTML sets the HTML template in memory and loads the card template via the utils
func SetPortfolioHTML(html string) {
	portfolioHTML = html
	if err := utils.LoadCardTemplate("components/card.html"); err != nil {
		panic("Failed to load card template: " + err.Error())
	}
}

// HomeHandler serves the HTML with dynamic cards
func HomeHandler(w http.ResponseWriter, r *http.Request) {
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

	cardsHTML := utils.BuildCardsHTML(articles)
	out := strings.Replace(portfolioHTML, "{{CARDS}}", cardsHTML, 1)
	w.Write([]byte(out))
}
