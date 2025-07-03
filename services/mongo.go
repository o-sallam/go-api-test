package services

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-api-test/models"
)

var mongoClient *mongo.Client
var articlesCollection *mongo.Collection

func InitMongo() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb+srv://Omar:123@cluster0.3ocieve.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}
	// Ping to verify connection and credentials
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("MongoDB ping error (check credentials): %v", err)
	}
	mongoClient = client
	articlesCollection = client.Database("Cluster0").Collection("articles")
}

func GetArticlesCollection() *mongo.Collection {
	return articlesCollection
}

func InsertArticle(article *models.Article) (*mongo.InsertOneResult, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return articlesCollection.InsertOne(ctx, article)
}

func EnsurePostsCollection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := mongoClient.Database("Cluster0")
	collections, err := db.ListCollectionNames(ctx, bson.M{"name": "posts"})
	if err != nil {
		return err
	}
	if len(collections) > 0 {
		return nil // Already exists
	}
	opts := options.CreateCollection()
	return db.CreateCollection(ctx, "posts", opts)
}

func AddPost(article *models.Article) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	db := mongoClient.Database("Cluster0")
	posts := db.Collection("posts")
	result, err := posts.InsertOne(ctx, article)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

func GetPostsCollection() *mongo.Collection {
	db := mongoClient.Database("Cluster0")
	return db.Collection("posts")
}
