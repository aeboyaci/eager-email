package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var database *mongo.Database

func Connect() {
	dbUri := "mongodb://aeboyaci:123456@127.0.0.1/?authSource=admin"

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(dbUri))

	if err != nil {
		log.Fatal("Cannot connect to MongoDb")
	}

	database = client.Database("eager_email")
}

func Disconnect() {
	database.Client().Disconnect(context.TODO())
}

func FindOne(collectionName string, filter interface{}) *mongo.SingleResult {
	collection := database.Collection(collectionName)

	result := collection.FindOne(context.TODO(), filter)

	return result
}

func FindMany(collectionName string, filter interface{}) ([]primitive.M, error) {
	collection := database.Collection(collectionName)

	cursor, err := collection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}

	var result []bson.M
	cursor.All(context.TODO(), &result)

	return result, nil
}

func InsertOne(collectionName string, document interface{}) (*mongo.InsertOneResult, error) {
	collection := database.Collection(collectionName)

	result, err := collection.InsertOne(context.TODO(), document)

	return result, err
}
