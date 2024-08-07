package model

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type Alert struct {
	Id          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	SearchType  string             `json:"search_type" bson:"search_type"`
	Location    string             `json:"location" bson:"location"`
	MaxPrice    int                `json:"max_price" bson:"max_price"`
	MinBedrooms int                `json:"min_bedrooms" bson:"min_bedrooms"`
	Subscribers []User             `json:"subscribers" bson:"subscribers"`
}

func (alert *Alert) Format() string {
	return alert.SearchType + " in " + alert.Location + " (max price of " + strconv.Itoa(alert.MaxPrice) + " and min " + strconv.Itoa(alert.MinBedrooms) + " rooms)"
}
