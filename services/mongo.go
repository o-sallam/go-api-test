package services

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"go-api-test/models"
)

var mongoClient *mongo.Client
var articlesCollection *mongo.Collection

func InitMongo() {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb+srv://Omar:<db_password>@cluster0.3ocieve.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
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
