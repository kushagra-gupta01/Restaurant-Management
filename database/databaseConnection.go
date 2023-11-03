package database

import(
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBinstance() *mongo.Client{
	MongoDb := "mongodb://localhost:27017"
	fmt.Print(MongoDb)
	
	ctx,cancel := context.WithTimeout(context.Background(),10*time.Second)
	client,err :=mongo.Connect(ctx,options.Client().ApplyURI(MongoDb))
	if err !=nil{
		log.Fatal(err)
	}
	defer cancel()
	fmt.Print("Connected to MongoDB")
	return client
}

var Client *mongo.Client = DBinstance()

func OpenCollection(client *mongo.Client) *mongo.Collection{
	var Collection *mongo.Collection = client.Database("restaurant").Collection("collectionName")
	return Collection
}	