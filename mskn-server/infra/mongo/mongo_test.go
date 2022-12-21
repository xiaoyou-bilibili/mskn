package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongo(t *testing.T) {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://192.168.1.10:8313"))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()
	coll := client.Database("demo").Collection("movies")
	coll.InsertOne(context.TODO(), map[string]interface{}{
		"age":  11,
		"name": "xiaoyou",
	})
}
