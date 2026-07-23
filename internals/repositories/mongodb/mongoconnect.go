package mongodb

import (
	"context"
	"fmt"
	"grpcapi/pkg/utils"
	"log"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func CreateMongoClient(ctx context.Context) (*mongo.Client, error) {
	client, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, utils.ErrorHandler(err, "Error connections to MongoDB")
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Println(err)
		return nil, utils.ErrorHandler(err, "Unable to ping database")
	}
	fmt.Println("Connected to MongoDB")
	return client, nil
}
