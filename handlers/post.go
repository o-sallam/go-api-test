package handlers

import (
	"context"
	"fmt"
	"go-api-test/models"
	"go-api-test/services"
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
	// Fill mainContent with article data (simple replace)
	mainContent = strings.ReplaceAll(mainContent, "{{title}}", article.Title)
	mainContent = strings.ReplaceAll(mainContent, "{{excerpt}}", article.Excerpt)
	mainContent = strings.ReplaceAll(mainContent, "{{content}}", article.Content)
	mainContent = strings.ReplaceAll(mainContent, "{{category}}", article.Category)
	mainContent = strings.ReplaceAll(mainContent, "{{coverImage}}", article.CoverImage)
	mainContent = strings.ReplaceAll(mainContent, "{{views}}", fmt.Sprintf("%d", article.Views))
	mainContent = strings.ReplaceAll(mainContent, "{{authorName}}", article.Author.Name)
	mainContent = strings.ReplaceAll(mainContent, "{{createdAt}}", article.CreatedAt)
	// Fill author aside
	authorAside := authorAsideTmpl
	authorAside = strings.ReplaceAll(authorAside, "{{authorImage}}", "/img/auth.webp")
	authorAside = strings.ReplaceAll(authorAside, "{{authorName}}", article.Author.Name)
	authorAside = strings.ReplaceAll(authorAside, "{{authorBio}}", "كاتب متخصص في عالم الطعام والمشروبات. يحب استكشاف الثقافات المختلفة من خلال مذاقاتها الفريدة.")
	authorAside = strings.ReplaceAll(authorAside, "{{authorJoin}}", "يناير 2020")
	authorAside = strings.ReplaceAll(authorAside, "{{authorArticles}}", "45")
	authorAside = strings.ReplaceAll(authorAside, "{{authorFollowers}}", "2.3k")
	authorAside = strings.ReplaceAll(authorAside, "{{articleStatsViews}}", fmt.Sprintf("%d", article.Views))
	authorAside = strings.ReplaceAll(authorAside, "{{articleComments}}", "23")
	authorAside = strings.ReplaceAll(authorAside, "{{articleStatsDate}}", article.CreatedAt)
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
