package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/utils"
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
	numDocs, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").CountDocuments(ctx, bson.D{{}})
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

func (dbClient *mongoDbClient) CreateProperties() map[string]string {
	property1 := model.Property{
		Url:               "https://url-test-1.com",
		Address:           "my keli1",
		Price:             1234,
		Type:              "Studio",
		NumDoubleBedrooms: 2,
		NumBathrooms:      1,
		Furnished:         true,
		LeaseType:         "minimum 1 lifetime",
		Description:       "this should be  very long description",
		ListingId:         "123456",
	}
	property2 := model.Property{
		Url:               "https://url-test-2.com",
		Address:           "my keli",
		Price:             1234,
		Type:              "Studio",
		NumDoubleBedrooms: 2,
		NumBathrooms:      1,
		Furnished:         true,
		LeaseType:         "minimum 1 lifetime",
		Description:       "this should be  very long description",
		ListingId:         "123457",
	}

	properties := []interface{}{
		property1,
		property2,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").InsertMany(ctx, properties)
	if res == nil {
		log.Fatalf(fmt.Sprintf("Error on CreateProperties() method: %v", err))
	}

	// ret := make(map[string]string)
	// for i, id := range res.InsertedIDs {
	// 	// TODO if needed.
	// }

	return map[string]string{
		"inserted": "ok",
	}
}

func (dbClient *mongoDbClient) DeleteProperties() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{}
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CountProperties() method: %v", err))
	}

	return map[string]int64{
		"count": res.DeletedCount,
	}
}

func (dbClient *mongoDbClient) FindPropertiesByListingIds() []model.Property {
	return make([]model.Property, 0)
}

func (dbClient *mongoDbClient) CreateAlertForUser(alert model.Alert, user model.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 1. Create User if missing
	dbClient.createUser(user)

	// 2. Upsert Alert with relevant User data.
	filter := bson.D{
		primitive.E{Key: "search_type", Value: alert.SearchType},
		primitive.E{Key: "location", Value: alert.Location},
		primitive.E{Key: "max_price", Value: alert.MaxPrice},
		primitive.E{Key: "min_bedrooms", Value: alert.MinBedrooms},
	}

	// $push (supports dupes) vs $addToSet (does not support dupes)
	update := bson.D{
		primitive.E{Key: "$addToSet", Value: bson.D{
			primitive.E{Key: "subscribers", Value: user}},
		},
	}
	opts := &options.UpdateOptions{
		Upsert: utils.BoolPtr(true),
	}
	_, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CreateAlertForUser()::2: %v", err))
	}

	return true
}

func (dbClient *mongoDbClient) ListAlertsForUser(user model.User) []model.Alert {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	//{"subscribers": {"$elemMatch": {"telegram_user_id": user.TelegramUserId}}}
	filter := bson.D{
		primitive.E{Key: "subscribers", Value: bson.D{
			primitive.E{Key: "$elemMatch", Value: bson.D{
				primitive.E{Key: "telegram_user_id", Value: user.TelegramUserId},
			}},
		}},
	}

	cur, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").Find(ctx, filter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on ListAlertsForUser()::1: %v", err))
	}

	res := make([]model.Alert, 4)
	err = cur.All(ctx, &res)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on ListAlertsForUser()::2: %v", err))
	}

	return res
}

func (dbClient *mongoDbClient) DeleteAlertForUser(alert model.Alert, user model.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 1. Remove User from Alert(s) matching specific criteria.
	filter := bson.D{
		primitive.E{Key: "search_type", Value: alert.SearchType},
		primitive.E{Key: "location", Value: alert.Location},
		primitive.E{Key: "max_price", Value: alert.MaxPrice},
		primitive.E{Key: "min_bedrooms", Value: alert.MinBedrooms},
	}

	update := bson.D{
		primitive.E{Key: "$pull", Value: bson.D{
			primitive.E{Key: "subscribers", Value: bson.D{
				primitive.E{Key: "telegram_user_id", Value: user.TelegramUserId},
			}}},
		},
	}
	_, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on DeleteAlertForUser()::1: %v", err))
	}

	// 2. Cleanup all alerts with no subscribers.
	filter = bson.D{
		primitive.E{Key: "subscribers", Value: bson.D{
			primitive.E{Key: "$exists", Value: true},
			primitive.E{Key: "$size", Value: 0},
		}},
	}
	_, err = dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").DeleteOne(ctx, filter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on DeleteAlertForUser()::2: %v", err))
	}

	return true
}

func (dbClient *mongoDbClient) createUser(user model.User) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "telegram_user_id", Value: user.TelegramUserId},
		primitive.E{Key: "telegram_chat_id", Value: user.TelegramChatId},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: filter},
	}
	opts := &options.UpdateOptions{
		Upsert: utils.BoolPtr(true),
	}

	_, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("users").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CreateUser(): %v", err))
	}
	return true
}

func (dbClient *mongoDbClient) DeleteAlerts() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{}
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on DeleteAlerts() method: %v", err))
	}

	return map[string]int64{
		"count": res.DeletedCount,
	}
}

func (dbClient *mongoDbClient) DeleteUsers() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{}
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("users").DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on DeleteUsers() method: %v", err))
	}

	return map[string]int64{
		"count": res.DeletedCount,
	}
}
