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

func getPrevNextArticles(slug string) (prev *models.Article, next *models.Article, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	posts := services.GetPostsCollection()
	// السابق: أول مقال slug أصغر من الحالي (حسب الترتيب الأبجدي)
	prevRes := posts.FindOne(ctx, map[string]interface{}{"slug": map[string]interface{}{"$lt": slug}})
	var prevArticle models.Article
	if err := prevRes.Decode(&prevArticle); err == nil {
		prev = &prevArticle
	}
	// التالي: أول مقال slug أكبر من الحالي (حسب الترتيب الأبجدي)
	nextRes := posts.FindOne(ctx, map[string]interface{}{"slug": map[string]interface{}{"$gt": slug}})
	var nextArticle models.Article
	if err := nextRes.Decode(&nextArticle); err == nil {
		next = &nextArticle
	}
	return prev, next, nil
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL: /{slug}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 2 || parts[1] == "" {
		utils.Show404(w)
		return
	}
	slug := parts[1]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Fetch article from DB
	var article models.Article
	err := services.GetPostsCollection().FindOne(ctx, map[string]interface{}{"slug": slug}).Decode(&article)
	if err != nil {
		utils.Show404(w)
		return
	}
	// جلب السابق والتالي
	prev, next, _ := getPrevNextArticles(slug)
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
		"prevLink":   "#",
		"prevTitle":  "لا يوجد مقال سابق",
		"prevImage":  "/img/blog",
		"nextLink":   "#",
		"nextTitle":  "لا يوجد مقال لاحق",
		"nextImage":  "/img/blog",
	}
	if prev != nil {
		mainFields["prevLink"] = "/" + prev.Slug
		mainFields["prevTitle"] = prev.Title
		mainFields["prevImage"] = prev.CoverImage
	}
	if next != nil {
		mainFields["nextLink"] = "/" + next.Slug
		mainFields["nextTitle"] = next.Title
		mainFields["nextImage"] = next.CoverImage
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

// PostPartialHTMLHandler returns only the <main>...</main> HTML of the post for a given slug
func PostPartialHTMLHandler(w http.ResponseWriter, r *http.Request) {
	// Extract slug from URL: /post-partial-html/{slug}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 3 || parts[2] == "" {
		utils.Show404(w)
		return
	}
	slug := parts[2]
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var article models.Article
	err := services.GetPostsCollection().FindOne(ctx, map[string]interface{}{"slug": slug}).Decode(&article)
	if err != nil {
		utils.Show404(w)
		return
	}
	postBytes, err := os.ReadFile("views/post.html")
	if err != nil {
		http.Error(w, "Post template not found", 500)
		return
	}
	postHTML := string(postBytes)
	mainStart := strings.Index(postHTML, "<main>")
	mainEnd := strings.Index(postHTML, "</main>")
	if mainStart == -1 || mainEnd == -1 {
		http.Error(w, "Invalid post template", 500)
		return
	}
	// extract only the content inside <main>...</main>
	mainInner := postHTML[mainStart+len("<main>") : mainEnd]
	mainFields := map[string]string{
		"title":      article.Title,
		"excerpt":    article.Excerpt,
		"content":    article.Content,
		"category":   article.Category,
		"coverImage": article.CoverImage,
		"views":      fmt.Sprintf("%d", article.Views),
		"authorName": article.Author.Name,
		"createdAt":  article.CreatedAt,
		"prevLink":   "#",
		"prevTitle":  "لا يوجد مقال سابق",
		"prevImage":  "/img/last-post",
		"nextLink":   "#",
		"nextTitle":  "لا يوجد مقال لاحق",
		"nextImage":  "/img/last-post",
	}
	prev, next, _ := getPrevNextArticles(slug)
	if prev != nil {
		mainFields["prevLink"] = "/" + prev.Slug
		mainFields["prevTitle"] = prev.Title
		mainFields["prevImage"] = prev.CoverImage
	}
	if next != nil {
		mainFields["nextLink"] = "/" + next.Slug
		mainFields["nextTitle"] = next.Title
		mainFields["nextImage"] = next.CoverImage
	}
	mainInner = utils.ReplacePlaceholders(mainInner, mainFields)
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
	mainInner = strings.Replace(mainInner, "{{AUTHOR_ASIDE}}", authorAside, 1)
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(mainInner))
}
