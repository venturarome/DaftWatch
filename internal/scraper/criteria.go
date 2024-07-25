package scraper

import (
	"errors"
	"fmt"
)

type Criteria struct {
	searchType   string
	location     string // TODO make this a list.
	propertyType string
	filters      []Filter
}

func (self *Criteria) buildQuery() (query string, err error) {
	if !self.isValid() { // TODO validate on object creation!
		return "", errors.New("Criteria is not valid") // TODO create custom error
	}
	// https://www.daft.ie/property-for-rent/dublin-city?rentalPrice_from=200&rentalPrice_to=1600&firstPublishDate_from=now-1d%2Fd

	query = fmt.Sprintf("https://www.daft.ie/%s/%s%s", self.mapSearchType(), self.mapLocation(), self.buildQueryParams())
	err = nil
	return
}

func (self *Criteria) isValid() (ok bool) {

	allowedFilters := map[string]string{
		// map values set to default value (different than "") if key is required!!
		"location":              "ireland", // TODO keep an eye on this one! There can be several locations. Maybe we dont want it. For simplicity, let's assume only 1.
		"rentalPrice_from":      "",
		"rentalPrice_to":        "",
		"firstPublishDate_from": "",
		// TODO fill it with all possible keys.
	}
	for _, param := range self.filters {
		_, found := allowedFilters[param.key]
		if !found {
			return false
		}
	}
	return true
}

// TODO create a DaftMapper.
func (self *Criteria) mapSearchType() string {
	allowedSearchTypes := map[string]string{
		// DaftWatch: Daft
		"buy":            "property-for-sale",
		"buy_commercial": "commercial-properties-for-sale",
		"buy_overseas":   "overseas-properties-for-sale",
		"buy_parking":    "parking-spaces-for-sale",
		"buy_new":        "new-homes-for-sale",

		"rent":            "property-for-rent",
		"rent_commercial": "commercial-properties-for-rent",
		"rent_overseas":   "overseas-properties-for-rent",
		"rent_parking":    "parking-spaces-for-rent",
		"rent_holidays":   "holiday-homes",
		"rent_student":    "student-accomodation-for-rent",

		"share":         "sharing",
		"share_student": "student-accomodation-for-share",

		// TODO fill all options
	}
	return allowedSearchTypes[self.searchType]
}

// TODO create a DaftMapper.
func (self *Criteria) mapLocation() string {
	allowedLocations := map[string]string{
		// DaftWatch: Daft
		"cork_county":     "cork",
		"cork_city":       "cork-city",
		"dublin_county":   "dublin",
		"dublin_city":     "dublin-city",
		"dublin_1":        "dublin-1-dublin",
		"dublin_2":        "dublin-2-dublin",
		"dublin_3":        "dublin-3-dublin",
		"dublin_4":        "dublin-4-dublin",
		"dublin_5":        "dublin-5-dublin",
		"dublin_6":        "dublin-6-dublin",
		"dublin_7":        "dublin-7-dublin",
		"dublin_8":        "dublin-8-dublin",
		"dublin_9":        "dublin-9-dublin",
		"dublin_10":       "dublin-10-dublin",
		"dublin_11":       "dublin-11-dublin",
		"dublin_12":       "dublin-12-dublin",
		"dublin_13":       "dublin-13-dublin",
		"dublin_14":       "dublin-14-dublin",
		"dublin_15":       "dublin-15-dublin",
		"dublin_16":       "dublin-16-dublin",
		"dublin_17":       "dublin-17-dublin",
		"dublin_18":       "dublin-18-dublin",
		"dublin_20":       "dublin-20-dublin",
		"dublin_22":       "dublin-22-dublin",
		"dublin_24":       "dublin-24-dublin",
		"galway_county":   "galway",
		"galway_city":     "galway-city",
		"limerick_county": "limerick",
		"limerick_city":   "limerick-city",
		// TODO fill all options
	}
	return allowedLocations[self.location]
}

// TODO create a DaftMapper.
func (criteria *Criteria) mapPropertyType() string {
	allowedPropertyTypes := map[string]string{
		// DaftWatch: Daft
		"houses":        "houses",
		"apartments":    "apartments",
		"sites":         "sites",
		"holiday_homes": "holiday-homes",
		// TODO fill all options
	}
	return allowedPropertyTypes[criteria.propertyType]
}

func (criteria *Criteria) buildQueryParams() string {
	if len(criteria.filters) == 0 {
		return ""
	}

	queryParams := "?"
	for _, filter := range criteria.filters {
		queryParams += fmt.Sprintf("%s=%s&", filter.key, filter.value)
	}
	return queryParams[:len(queryParams)-1]
}
