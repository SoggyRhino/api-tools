/*
This file is responsible for providing various useful database functions.
*/
package uploader

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/UTDNebula/api-tools/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Replace instances and check
type DBSingleton struct {
	client *mongo.Client
}

var dbInstance *DBSingleton
var once sync.Once

func connectDB() *mongo.Client {
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		opts := options.Client().ApplyURI(getEnvMongoURI())

		client, err := mongo.Connect(ctx, opts)
		if err != nil {
			log.Fatalf("Unable to create MongoDB client and connect to database")
		}

		// ping the database
		err = client.Ping(ctx, nil)
		if err != nil {
			log.Fatalf("Unable to ping database")
		}

		log.Println("Connected to MongoDB")

		dbInstance = &DBSingleton{
			client: client,
		}
	})

	return dbInstance.client
}

func getCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	collection := client.Database("combinedDB").Collection(collectionName)
	return collection
}

func getEnvMongoURI() string {
	uri, err := utils.GetEnv("MONGODB_URI")
	if err != nil {
		panic(err)
	}
	return uri
}
