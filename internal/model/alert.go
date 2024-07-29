package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Alert struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	SearchType  string             `json:"search_type" bson:"search_type"`
	Location    string             `json:"location" bson:"location"`
	MaxPrice    uint16             `json:"max_price" bson:"max_price"`
	MinBedrooms uint16             `json:"min_bedrooms" bson:"min_bedrooms"`
	Subscribers []User             `json:"subscribers" bson:"subscribers"`
}
