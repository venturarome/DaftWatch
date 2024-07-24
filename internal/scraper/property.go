package scraper

type Property struct {
	url               string // TODO consider using net.url package
	address           string
	price             uint16
	propertyType      string // TODO create PropertyType enum or the like
	numDoubleBedrooms int
	numSingleBedrooms int
	numBathrooms      int
	availableFrom     string // TODO consider using date package
	isFurnished       bool
	leaseType         string // TODO create LeaseType enum or the like
	description       string
	advertiser        string // TODO create Advertiser struct {name string, phones []string, }
}
