package models

type Article struct {
	ID         string   `bson:"_id,omitempty" json:"id"`
	Slug       string   `bson:"slug" json:"slug"`
	Title      string   `bson:"title" json:"title"`
	Excerpt    string   `bson:"excerpt" json:"excerpt"`
	Content    string   `bson:"content" json:"content"`
	Category   string   `bson:"category" json:"category"`
	Tags       []string `bson:"tags" json:"tags"`
	Author     Author   `bson:"author" json:"author"`
	CoverImage string   `bson:"coverImage" json:"coverImage"`
	Views      int      `bson:"views" json:"views"`
	CreatedAt  string   `bson:"createdAt" json:"createdAt"`
	UpdatedAt  string   `bson:"updatedAt" json:"updatedAt"`
	Published  bool     `bson:"published" json:"published"`
}

type Author struct {
	ID   string `bson:"id" json:"id"`
	Name string `bson:"name" json:"name"`
}
