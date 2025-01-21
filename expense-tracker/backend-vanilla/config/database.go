package config

import (
    "context"
    "log"
    "os"
    "time"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectDB() {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Create client
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("MONGODB_URI")))
    if err != nil {
        log.Fatal(err)
    }

    // Ping database
    if err := client.Ping(ctx, nil); err != nil {
        log.Fatal(err)
    }

    DB = client.Database(os.Getenv("DB_NAME"))
    log.Printf("Connected to MongoDB: %s\n", os.Getenv("DB_NAME"))

    // Create indexes
    createIndexes(ctx)
}

func createIndexes(ctx context.Context) {
    // User indexes
    userColl := DB.Collection("users")
    _, err := userColl.Indexes().CreateOne(ctx, mongo.IndexModel{
        Keys:    bson.D{{Key: "email", Value: 1}},
        Options: options.Index().SetUnique(true),
    })
    if err != nil {
        log.Fatal(err)
    }

    // Expense indexes
    expenseColl := DB.Collection("expenses")
    _, err = expenseColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
        {
            Keys: bson.D{
                {Key: "userId", Value: 1},
                {Key: "date", Value: -1},
            },
        },
        {
            Keys: bson.D{{Key: "category", Value: 1}},
        },
    })
    if err != nil {
        log.Fatal(err)
    }
}
