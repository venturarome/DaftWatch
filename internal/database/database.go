package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Database interface {
	Health() map[string]string
	CreateProperties()
	CreateAlert()
	// TODO add here all methods to interact with databases. Will be implemented by all DB clientes (so far, only MongoDB)
}

type mongoDb struct {
	db *mongo.Client
}

// from go-blueprint. Probable is deprecated.
// var (
// 	host = os.Getenv("DB_HOST")
// 	port = os.Getenv("DB_PORT")
// 	database = os.Getenv("DB_DATABASE")
// 	opts = options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%s", host, port))
// )

func InstanceMongoDb() Database {
	dbUri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s.zcallcn.mongodb.net/?retryWrites=true&w=majority&appName=%s",
		os.Getenv("MONGO_DB_USERNAME"),
		os.Getenv("MONGO_DB_PASSWORD"),
		os.Getenv("MONGO_DB_CLUSTER"),
		os.Getenv("MONGO_DB_CLUSTER"),
	)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)

	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		log.Fatal(err)

	}
	return &mongoDb{
		db: client,
	}
}

func (conn *mongoDb) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := conn.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (conn *mongoDb) CreateProperties() {
	// TODO

	// Other example:
	// // Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// 	}
	// 	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
}

func (conn *mongoDb) CreateAlert() {
	// TODO
}
