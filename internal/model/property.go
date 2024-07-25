package model

type Property struct {
	Id                string `json:"_id" bson:"_id"`
	Url               string `json:"url" bson:"url"` // TODO consider using net.url package
	Address           string `json:"address" bson:"address"`
	Price             uint16 `json:"price" bson:"price"`
	Type              string `json:"type" bson:"type"` // TODO create PropertyType enum or the like
	NumDoubleBedrooms uint16 `json:"num_double_bedrooms" bson:"num_double_bedrooms"`
	NumSingleBedrooms uint16 `json:"num_single_bedrooms" bson:"num_single_bedrooms"`
	NumBathrooms      uint16 `json:"num_bathrooms" bson:"num_bathrooms"`
	AvailableFrom     string `json:"available_from" bson:"available_from"` // TODO consider using date package
	Furnished         bool   `json:"furnished" bson:"furnished"`
	LeaseType         string `json:"lease_type" bson:"lease_type"` // TODO create LeaseType enum or the like
	Description       string `json:"description" bson:"description"`
	ListingId         string `json:"listing_id" bson:"listing_id"`
}
