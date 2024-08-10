package scraper

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/venturarome/DaftWatch/internal/model"
)

type Criteria struct {
	SearchType string
	Location   string // TODO make this a list.
	//PropertyType string
	Filters []Filter
}

var SearchTypesMap = map[string]string{
	// DaftWatch: Daft
	"Buy": "property-for-sale",
	// "buy_commercial": "commercial-properties-for-sale",
	// "buy_overseas":   "overseas-properties-for-sale",
	// "buy_parking":    "parking-spaces-for-sale",
	// "buy_new":        "new-homes-for-sale",

	"Rent": "property-for-rent",
	// "rent_commercial": "commercial-properties-for-rent",
	// "rent_overseas":   "overseas-properties-for-rent",
	// "rent_parking":    "parking-spaces-for-rent",
	// "rent_holidays":   "holiday-homes",
	// "rent_student":    "student-accomodation-for-rent",

	// "share":         "sharing",
	// "share_student": "student-accomodation-for-share",

	// TODO fill all options
}

var LocationsMap = map[string]string{
	// DaftWatch: Daft

	//"belfast": "belfast-city",
	// "cork_county": "cork",
	"Cork": "cork-city",
	///"dublin_county": "dublin",
	"Dublin":       "dublin-city",
	"Dublin North": "north-dublin-city-dublin",
	"Dublin South": "south-dublin-city-dublin",
	"Dublin 01":    "dublin-1-dublin",
	"Dublin 02":    "dublin-2-dublin",
	"Dublin 03":    "dublin-3-dublin",
	"Dublin 04":    "dublin-4-dublin",
	"Dublin 05":    "dublin-5-dublin",
	"Dublin 06":    "dublin-6-dublin",
	"Dublin 07":    "dublin-7-dublin",
	"Dublin 08":    "dublin-8-dublin",
	"Dublin 09":    "dublin-9-dublin",
	"Dublin 10":    "dublin-10-dublin",
	"Dublin 11":    "dublin-11-dublin",
	"Dublin 12":    "dublin-12-dublin",
	"Dublin 13":    "dublin-13-dublin",
	"Dublin 14":    "dublin-14-dublin",
	"Dublin 15":    "dublin-15-dublin",
	"Dublin 16":    "dublin-16-dublin",
	"Dublin 17":    "dublin-17-dublin",
	"Dublin 18":    "dublin-18-dublin",
	"Dublin 20":    "dublin-20-dublin",
	"Dublin 22":    "dublin-22-dublin",
	// "galway_county":   "galway",
	"Galway": "galway-city",
	// "limerick_county": "limerick",
	"Limerick": "limerick-city",
}

// TODO move to "utils"? Create a specific Writer?
func CreateCriteriaFromAlert(alert model.Alert) Criteria {
	return Criteria{
		SearchType: alert.SearchType,
		Location:   alert.Location,
		Filters: []Filter{
			{Key: "maxPrice", Value: strconv.Itoa(alert.MaxPrice)},
			{Key: "minBedrooms", Value: strconv.Itoa(alert.MinBedrooms)},
			{Key: "firstPosted", Value: "now-20m"}, // We force to only check very recent listings (last 20 mins), as only want properties from now on.
		},
	}
}

func (criteria *Criteria) buildQuery(baseUrl string) (string, error) {
	if !criteria.isValid() { // TODO validate on object creation?
		return "", errors.New("Criteria is not valid") // TODO create custom error
	}

	query := fmt.Sprintf(
		"%s/%s/%s%s",
		baseUrl,
		criteria.mapSearchType(),
		criteria.mapLocation(),
		criteria.buildQueryParams(),
	)

	return query, nil
}

func (criteria *Criteria) isValid() (ok bool) {
	var valid bool

	_, valid = SearchTypesMap[criteria.SearchType]
	if !valid {
		return false
	}

	_, valid = LocationsMap[criteria.Location]
	if !valid {
		return false
	}

	for _, filter := range criteria.Filters {
		if !filter.isValid() {
			return false
		}
	}

	return true
}

// TODO create a DaftMapper.
func (criteria *Criteria) mapSearchType() string {
	return SearchTypesMap[criteria.SearchType]
}

// TODO create a DaftMapper.
func (criteria *Criteria) mapLocation() string {
	return LocationsMap[criteria.Location]
}

// // TODO create a DaftMapper.
// func (criteria *Criteria) mapPropertyType() string {
// 	allowedPropertyTypes := map[string]string{
// 		// DaftWatch: Daft
// 		"houses":        "houses",
// 		"apartments":    "apartments",
// 		"sites":         "sites",
// 		"holiday_homes": "holiday-homes",
// 		// TODO fill all options
// 	}
// 	return allowedPropertyTypes[criteria.propertyType]
// }

func (criteria *Criteria) buildQueryParams() string {
	if len(criteria.Filters) == 0 {
		return ""
	}

	queryParams := "?"
	for _, filter := range criteria.Filters {
		queryParams += fmt.Sprintf(
			"%s&",
			filter.toQueryParamForSearchType(criteria.SearchType),
		)

	}
	return queryParams[:len(queryParams)-1]
}
