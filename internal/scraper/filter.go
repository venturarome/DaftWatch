package scraper

import "fmt"

type Filter struct {
	Key   string
	Value string
}

var allowedFilters = map[string]string{
	"maxPrice":    "9999",
	"minBedrooms": "0",
}

var filtersMap = map[string]map[string]string{
	// DaftWatch: Daft
	"rent": {
		"maxPrice":    "rentalPrice_to",
		"minBedrooms": "numBeds_from",
	},
	"buy": {
		"maxPrice":    "salePrice_to",
		"minBedrooms": "numBeds_from",
	},
	// TODO fill it with all possible keys.
}

func (filter Filter) isValid() bool {
	_, valid := allowedFilters[filter.Key]
	return valid
}

func (filter Filter) toQueryParamForSearchType(searchType string) string {
	mappedKey := filtersMap[searchType][filter.Key]
	return fmt.Sprintf("%s=%s", mappedKey, filter.Value)
}
