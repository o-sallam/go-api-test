package handlers

import (
	"go-api-test/models"
	"go-api-test/utils"
	"net/http"
	"os"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/html"
)

var (
	portfolioBody string
	headerHTML    string
	footerHTML    string
	layoutHTML    string
)

// SetPortfolioHTML sets the HTML template in memory and loads the card template via the utils
func SetPortfolioHTML(body string) {
	portfolioBody = body
	if err := utils.LoadCardTemplate("components/card.html"); err != nil {
		panic("Failed to load card template: " + err.Error())
	}
	h, err := os.ReadFile("components/header.html")
	if err != nil {
		panic("Failed to load header.html: " + err.Error())
	}
	headerHTML = string(h)
	f, err := os.ReadFile("components/footer.html")
	if err != nil {
		panic("Failed to load footer.html: " + err.Error())
	}
	footerHTML = string(f)
	l, err := os.ReadFile("components/layout.html")
	if err != nil {
		panic("Failed to load layout.html: " + err.Error())
	}
	layoutHTML = string(l)
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
	body := strings.Replace(portfolioBody, "{{CARDS}}", cardsHTML, 1)
	out := layoutHTML
	out = strings.Replace(out, "{{HEADER}}", headerHTML, 1)
	out = strings.Replace(out, "{{BODY}}", body, 1)
	out = strings.Replace(out, "{{FOOTER}}", footerHTML, 1)

	m := minify.New()
	m.AddFunc("text/html", html.Minify)
	minified, err := m.String("text/html", out)
	if err != nil {
		w.Write([]byte(out)) // fallback to unminified
		return
	}
	w.Write([]byte(minified))
}
