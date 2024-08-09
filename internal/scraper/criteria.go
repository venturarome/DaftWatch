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

var searchTypesMap = map[string]string{
	// DaftWatch: Daft
	"buy": "property-for-sale",
	// "buy_commercial": "commercial-properties-for-sale",
	// "buy_overseas":   "overseas-properties-for-sale",
	// "buy_parking":    "parking-spaces-for-sale",
	// "buy_new":        "new-homes-for-sale",

	"rent": "property-for-rent",
	// "rent_commercial": "commercial-properties-for-rent",
	// "rent_overseas":   "overseas-properties-for-rent",
	// "rent_parking":    "parking-spaces-for-rent",
	// "rent_holidays":   "holiday-homes",
	// "rent_student":    "student-accomodation-for-rent",

	// "share":         "sharing",
	// "share_student": "student-accomodation-for-share",

	// TODO fill all options
}

var locationsMap = map[string]string{
	// DaftWatch: Daft

	// Simplified version
	"cork":     "cork-city",
	"dublin":   "dublin-city",
	"galway":   "galway-city",
	"limerick": "limerick-city",

	// Full version
	// "cork_county":     "cork",
	// "cork_city":       "cork-city",
	// "dublin_county":   "dublin",
	// "dublin_city":     "dublin-city",
	// "dublin_1":        "dublin-1-dublin",
	// "dublin_2":        "dublin-2-dublin",
	// "dublin_3":        "dublin-3-dublin",
	// "dublin_4":        "dublin-4-dublin",
	// "dublin_5":        "dublin-5-dublin",
	// "dublin_6":        "dublin-6-dublin",
	// "dublin_7":        "dublin-7-dublin",
	// "dublin_8":        "dublin-8-dublin",
	// "dublin_9":        "dublin-9-dublin",
	// "dublin_10":       "dublin-10-dublin",
	// "dublin_11":       "dublin-11-dublin",
	// "dublin_12":       "dublin-12-dublin",
	// "dublin_13":       "dublin-13-dublin",
	// "dublin_14":       "dublin-14-dublin",
	// "dublin_15":       "dublin-15-dublin",
	// "dublin_16":       "dublin-16-dublin",
	// "dublin_17":       "dublin-17-dublin",
	// "dublin_18":       "dublin-18-dublin",
	// "dublin_20":       "dublin-20-dublin",
	// "dublin_22":       "dublin-22-dublin",
	// "dublin_24":       "dublin-24-dublin",
	// "galway_county":   "galway",
	// "galway_city":     "galway-city",
	// "limerick_county": "limerick",
	// "limerick_city":   "limerick-city",
	// TODO fill all options
}

// TODO move to "utils"? Create a specific Writer?
func CreateCriteriaFromAlert(alert model.Alert) Criteria {
	return Criteria{
		SearchType: alert.SearchType,
		Location:   alert.Location,
		Filters: []Filter{
			{Key: "maxPrice", Value: strconv.Itoa(alert.MaxPrice)},
			{Key: "minBedrooms", Value: strconv.Itoa(alert.MinBedrooms)},
			{Key: "firstPosted", Value: "now-20m"},
		},
	}
}

func (criteria *Criteria) buildQuery(baseUrl string) (string, error) {
	if !criteria.isValid() { // TODO validate on object creation?
		return "", errors.New("Criteria is not valid") // TODO create custom error
	}
	// TODO add Filter with firstPublishDate_from. Maybe try "now-1h".
	// [https://www.daft.ie]/property-for-rent/dublin-city?rentalPrice_to=1600&firstPublishDate_from=now-1d

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

	_, valid = searchTypesMap[criteria.SearchType]
	if !valid {
		return false
	}

	_, valid = locationsMap[criteria.Location]
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
	return searchTypesMap[criteria.SearchType]
}

// TODO create a DaftMapper.
func (criteria *Criteria) mapLocation() string {
	return locationsMap[criteria.Location]
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
