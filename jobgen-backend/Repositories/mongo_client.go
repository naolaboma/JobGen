package repositories

import (
    "context"
    "jobgen-backend/Infrastructure"
    "log"
    "time"

    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoClient() *mongo.Client {
    client, err := mongo.NewClient(options.Client().ApplyURI(infrastructure.Env.MongoDBURI))
    if err != nil {
        log.Fatal("Failed to create MongoDB client:", err)
    }
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := client.Connect(ctx); err != nil {
        log.Fatal("Failed to connect to MongoDB:", err)
    }
    return client
}

func GetDatabase(client *mongo.Client) *mongo.Database {
    return client.Database(infrastructure.Env.DBName)
}
