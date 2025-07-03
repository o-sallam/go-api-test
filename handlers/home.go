package handlers

import (
	"context"
	"fmt"
	"go-api-test/models"
	"go-api-test/services"
	"go-api-test/utils"
	"net/http"
	"os"
	"strings"
	"time"

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

	// Fetch articles from the database
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cur, err := services.GetPostsCollection().Find(ctx, map[string]interface{}{})
	if err != nil {
		http.Error(w, "Failed to fetch articles", 500)
		return
	}
	var articles []models.Article
	if err := cur.All(ctx, &articles); err != nil {
		http.Error(w, "Failed to decode articles", 500)
		return
	}
	// Convert to DummyArticle for card rendering
	var dummyArticles []models.DummyArticle
	for _, a := range articles {
		// Format date as short date (YYYY-MM-DD)
		shortDate := a.CreatedAt
		if len(shortDate) >= 10 {
			shortDate = shortDate[:10]
		}
		dummyArticles = append(dummyArticles, models.DummyArticle{
			ALT:      a.Title + " - صورة توضيحية",
			IMG:      a.CoverImage,
			CATEGORY: a.Category,
			LINK:     "/" + a.Slug,
			TITLE:    a.Title,
			EXCERPT:  a.Excerpt,
			VIEWS:    fmt.Sprintf("%d", a.Views),
			AUTHOR:   a.Author.Name,
			DATE:     shortDate,
			Slug:     a.Slug,
		})
	}
	cardsHTML := utils.BuildCardsHTML(dummyArticles)
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
