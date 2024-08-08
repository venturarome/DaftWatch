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

// GO and BSON: https://www.mongodb.com/docs/drivers/go/current/fundamentals/bson/#data-types

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

func (dbClient *mongoDbClient) CreateUser(user model.User) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "telegram_user_id", Value: user.TelegramUserId},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: filter},
	}
	opts := &options.UpdateOptions{
		Upsert: utils.BoolPtr(true),
	}

	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("users").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CreateUser(): %v", err))
	}

	return map[string]interface{}{
		"MatchedCount":  res.MatchedCount,
		"ModifiedCount": res.ModifiedCount,
		"UpsertedCount": res.UpsertedCount,
		"UpsertedID":    res.UpsertedID,
	}
}

func (dbClient *mongoDbClient) ListAlertsForUser(user model.User) []model.Alert {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

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

func (dbClient *mongoDbClient) AddSubscriberToAlert(alert model.Alert, user model.User) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "search_type", Value: alert.SearchType},
		primitive.E{Key: "location", Value: alert.Location},
		primitive.E{Key: "max_price", Value: alert.MaxPrice},
		primitive.E{Key: "min_bedrooms", Value: alert.MinBedrooms},
	}

	// $push (supports dupes) vs $addToSet (does not support dupes)
	update := bson.D{
		primitive.E{Key: "$addToSet", Value: bson.D{
			primitive.E{Key: "subscribers", Value: user.TelegramUserId}},
		},
	}
	opts := &options.UpdateOptions{
		Upsert: utils.BoolPtr(true),
	}
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on AddSubscriberToAlert(): %v", err))
	}

	return map[string]interface{}{
		"MatchedCount":  res.MatchedCount,
		"ModifiedCount": res.ModifiedCount,
		"UpsertedCount": res.UpsertedCount,
		"UpsertedID":    res.UpsertedID,
	}
}

func (dbClient *mongoDbClient) RemoveSubscriberFromAlert(alert model.Alert, user model.User) bool {
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
		log.Fatalf(fmt.Sprintf("Error on RemoveSubscriberFromAlert()::1: %v", err))
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
		log.Fatalf(fmt.Sprintf("Error on RemoveSubscriberFromAlert()::2: %v", err))
	}

	return true
}

func (DbClient *mongoDbClient) SetPropertiesToAlert(alert model.Alert, properties []model.Property) map[string]interface{} {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "_id", Value: alert.Id},
	}

	propertyListingIds := utils.MapSlice(
		properties,
		func(p model.Property) string {
			return p.ListingId
		},
	)
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "properties", Value: propertyListingIds},
		}},
	}

	res, err := DbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("alerts").UpdateOne(ctx, filter, update)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on SetPropertiesToAlert(): %v", err))
	}

	return map[string]interface{}{
		"MatchedCount":  res.MatchedCount,
		"ModifiedCount": res.ModifiedCount,
		"UpsertedCount": res.UpsertedCount,
		"UpsertedID":    res.UpsertedID,
	}
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

func (dbClient *mongoDbClient) CreateProperty(property model.Property) map[string]interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{
		primitive.E{Key: "listing_id", Value: property.ListingId},
	}
	update := bson.D{
		primitive.E{Key: "$set", Value: bson.D{
			primitive.E{Key: "url", Value: property.Url},
			primitive.E{Key: "address", Value: property.Address},
			primitive.E{Key: "price", Value: property.Price},
			primitive.E{Key: "type", Value: property.Type},
			primitive.E{Key: "num_double_bedrooms", Value: property.NumDoubleBedrooms},
			primitive.E{Key: "num_single_bedrooms", Value: property.NumSingleBedrooms},
			primitive.E{Key: "num_bathrooms", Value: property.NumBathrooms},
			primitive.E{Key: "avaliable_from", Value: property.AvailableFrom},
			primitive.E{Key: "furnished", Value: property.Furnished},
			primitive.E{Key: "lease_type", Value: property.LeaseType},
			primitive.E{Key: "description", Value: property.Description},
			primitive.E{Key: "listing_id", Value: property.ListingId},
		}},
	}
	opts := &options.UpdateOptions{
		Upsert: utils.BoolPtr(true),
	}

	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on CreateProperty(): %v", err))
	}

	return map[string]interface{}{
		"MatchedCount":  res.MatchedCount,
		"ModifiedCount": res.ModifiedCount,
		"UpsertedCount": res.UpsertedCount,
		"UpsertedID":    res.UpsertedID,
	}
}

func (dbClient *mongoDbClient) CreateProperties(properties []model.Property) map[string]interface{} {

	var matchedCount, modifiedCount, upsertedCount int64 = 0, 0, 0
	var upsertedIDs []primitive.ObjectID = make([]primitive.ObjectID, 0, len(properties))
	for _, property := range properties {
		res := dbClient.CreateProperty(property)
		matchedCount += res["MatchedCount"].(int64)
		modifiedCount += res["ModifiedCount"].(int64)
		upsertedCount += res["UpsertedCount"].(int64)
		if v, ok := res["UpsertedID"].(primitive.ObjectID); ok {
			upsertedIDs = append(upsertedIDs, v)
		}

	}

	return map[string]interface{}{
		"MatchedCount":  matchedCount,
		"ModifiedCount": modifiedCount,
		"UpsertedCount": upsertedCount,
		"UpsertedIDs":   upsertedIDs,
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

func (dbClient *mongoDbClient) DeleteProperties() map[string]int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	filter := bson.D{}
	res, err := dbClient.db.Database(os.Getenv("MONGO_DB_DATABASE")).Collection("properties").DeleteMany(ctx, filter)
	if err != nil {
		log.Fatalf(fmt.Sprintf("Error on DeleteProperties() method: %v", err))
	}

	return map[string]int64{
		"count": res.DeletedCount,
	}
}
