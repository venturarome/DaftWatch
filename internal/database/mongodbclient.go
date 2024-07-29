package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/venturarome/DaftWatch/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoDbClient struct {
	db *mongo.Client
}

func InstanceMongoDb() DbClient {
	dbUri := fmt.Sprintf(
		"mongodb+srv://%s:%s@%s.zcallcn.mongodb.net/?retryWrites=true&w=majority&appName=%s",
		os.Getenv("MONGO_DB_USERNAME"),
		os.Getenv("MONGO_DB_PASSWORD"),
		os.Getenv("MONGO_DB_CLUSTER"),
		os.Getenv("MONGO_DB_CLUSTER"),
	)
	serverApiOptions := options.ServerAPI(options.ServerAPIVersion1)
	bsonOptions := &options.BSONOptions{
		UseJSONStructTags: true,
		NilSliceAsEmpty:   true,
	}

	opts := options.Client().ApplyURI(dbUri).SetServerAPIOptions(serverApiOptions).SetBSONOptions(bsonOptions)

	client, err := mongo.Connect(context.Background(), opts)

	if err != nil {
		log.Fatal(err)
	}

	return &mongoDbClient{
		db: client,
	}
}

func (dbClient *mongoDbClient) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := dbClient.db.Ping(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("db down: %v", err))
	}

	return map[string]string{
		"message": "It's healthy",
	}
}

func (dbClient *mongoDbClient) CountProperties() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	numDocs, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").CountDocuments(ctx, nil)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CountProperties() method: %v", err))
	}

	return map[string]int64{
		"count": numDocs,
	}
}

// GO and BSON: https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/#data-types

func (dbClient *mongoDbClient) CreateProperty() map[string]string {
	property := model.Property{
		Url:               "https://url-test.com",
		Address:           "my keli",
		Price:             1234,
		Type:              "Studio",
		NumDoubleBedrooms: 2,
		NumBathrooms:      1,
		Furnished:         true,
		LeaseType:         "minimum 1 lifetime",
		Description:       "this should be  very long description",
		ListingId:         "123456",
		// Omitted NumSingleBedrooms
		// Omitted AvailableFrom
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").InsertOne(ctx, property)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CountProperties() method: %v", err))
	}

	return map[string]string{
		"insertedId": res.InsertedID.(primitive.ObjectID).String(),
	}
}

func (dbClient *mongoDbClient) DeleteProperties() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").DeleteMany(ctx, bson.D{})
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CountProperties() method: %v", err))
	}

	return map[string]int64{
		"count": res.DeletedCount,
	}
}

func (dbClient *mongoDbClient) CreateProperties() map[string]string {
	// TODO

	// Other example:
	// // Send a ping to confirm a successful connection
	// if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
	// 	panic(err)
	// 	}
	// 	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")
	return map[string]string{
		"TODO": "TODO",
	}
}

func (dbClient *mongoDbClient) CreateAlert() map[string]string {
	return map[string]string{
		"TODO": "TODO",
	}
}
