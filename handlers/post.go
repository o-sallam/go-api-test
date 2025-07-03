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
)

var (
	postLayoutHTML  string
	authorAsideTmpl string
)

func InitPostTemplates() {
	// Load layout and author aside
	layoutBytes, err := os.ReadFile("components/layout.html")
	if err != nil {
		panic("Failed to load layout.html: " + err.Error())
	}
	// Replace main.css with post.css for post layout
	postLayoutHTML = strings.Replace(string(layoutBytes), "/css/main.css", "/css/post.css", 1)
	a, err := os.ReadFile("components/author-aside.html")
	if err != nil {
		panic("Failed to load author-aside.html: " + err.Error())
	}
	authorAsideTmpl = string(a)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL: /{slug}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 || parts[1] == "" {
		http.NotFound(w, r)
		return
	}
	slug := parts[1]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Fetch article from DB
	var article models.Article
	err := services.GetPostsCollection().FindOne(ctx, map[string]interface{}{"slug": slug}).Decode(&article)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	// Load post.html (main content template)
	postBytes, err := os.ReadFile("views/post.html")
	if err != nil {
		http.Error(w, "Post template not found", 500)
		return
	}
	postHTML := string(postBytes)
	// Remove repeated html/head/body from post.html, keep only <main>...</main>
	mainStart := strings.Index(postHTML, "<main>")
	mainEnd := strings.Index(postHTML, "</main>")
	if mainStart == -1 || mainEnd == -1 {
		http.Error(w, "Invalid post template", 500)
		return
	}
	mainContent := postHTML[mainStart : mainEnd+len("</main>")]
	// Fill mainContent with article data (dynamic replace)
	mainFields := map[string]string{
		"title":      article.Title,
		"excerpt":    article.Excerpt,
		"content":    article.Content,
		"category":   article.Category,
		"coverImage": article.CoverImage,
		"views":      fmt.Sprintf("%d", article.Views),
		"authorName": article.Author.Name,
		"createdAt":  article.CreatedAt,
	}
	mainContent = utils.ReplacePlaceholders(mainContent, mainFields)
	// Fill author aside (dynamic replace)
	authorAside := authorAsideTmpl
	authorFields := map[string]string{
		"authorImage":       "/img/auth.webp",
		"authorName":        article.Author.Name,
		"authorBio":         "كاتب متخصص في عالم الطعام والمشروبات. يحب استكشاف الثقافات المختلفة من خلال مذاقاتها الفريدة.",
		"authorJoin":        "يناير 2020",
		"authorArticles":    "45",
		"authorFollowers":   "2.3k",
		"articleStatsViews": fmt.Sprintf("%d", article.Views),
		"articleComments":   "23",
		"articleStatsDate":  article.CreatedAt,
	}
	authorAside = utils.ReplacePlaceholders(authorAside, authorFields)
	// Insert authorAside into mainContent at the placeholder
	mainContent = strings.Replace(mainContent, "{{AUTHOR_ASIDE}}", authorAside, 1)
	// Compose final HTML using layout
	out := postLayoutHTML
	out = strings.Replace(out, "{{HEADER}}", headerHTML, 1)
	out = strings.Replace(out, "{{BODY}}", mainContent, 1)
	out = strings.Replace(out, "{{FOOTER}}", footerHTML, 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(out))
}
