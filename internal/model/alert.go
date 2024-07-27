package model

type Alert struct {
	Id          string `json:"_id" bson:"_id"`
	SearchType  string `json:"search_type" bson:"search_type"`
	Location    string `json:"location" bson:"location"`
	MaxPrice    uint16 `json:"max_price" bson:"max_price"`
	MinBedrooms uint16 `json:"min_bedrooms" bson:"min_bedrooms"`
	Subscribers []User `json:"subscribers" bson:"subscribers"`
}
