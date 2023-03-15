package services

import (
	"context"
	"fmt"
	"log"
	"os"

	AppConstant "github.com/HousewareHQ/backend-engineering-octernship/api/server/constants"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func DBClient() *mongo.Client { //RETURNS MONGODB client instance
	err := godotenv.Load("../../.env") //Load .env file to access db credentials
	if err != nil {
		log.Fatal("Failed to load ", err.Error())
	}

	MongoDBCreds := os.Getenv("MONGODB_CREDURL") //Store url in a variable
	// Use the SetServerAPIOptions() method to set the Stable API version to 1
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(MongoDBCreds).SetServerAPIOptions(serverAPI)
	// Create a new client and connect to the server
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatal("Failed to Connect to mongoDB")
	}

	// defer func() {
	// 	//At the end of this function , disconnect from mongoDB
	// 	if err = client.Disconnect(context.TODO()); err != nil {
	// 		panic(err)
	// 	}
	// }()

	fmt.Println("Connected to database..")

	return client

}

// Global variable to MongoDB client instance
var Client *mongo.Client = DBClient()

// Returns MongoDB collection
func OpenCollection(client *mongo.Client, collecName string) *mongo.Collection {
	var Collection *mongo.Collection = (*mongo.Collection)(client.Database(AppConstant.DB_NAME).Collection(collecName))
	return Collection
}

// For test purpose
func PingDB() {
	err := Client.Ping(context.TODO(), readpref.Primary())
	if err != nil {
		fmt.Println("Not connected to DB..")
	}

}
