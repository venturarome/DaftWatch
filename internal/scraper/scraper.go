package scraper

import (
	"fmt"
	"path"
	"strings"

	"github.com/gocolly/colly"
	"github.com/venturarome/DaftWatch/internal/model"
	"github.com/venturarome/DaftWatch/internal/utils"
)

const HTTP_SECURE string = "https://"
const BASE_URL_DAFT string = "www.daft.ie"

func Scrape(criteria Criteria) map[string]model.Property {
	fmt.Println("Entering scraper::Scrape")

	// Map to store scraped data
	propertyIdUrlMap := make(map[string]model.Property)

	// Main page collector
	mainCollector := colly.NewCollector(
		colly.AllowedDomains(BASE_URL_DAFT),
		colly.DetectCharset(),
	)

	// Property page collector
	propertyCollector := mainCollector.Clone()

	// Define behaviour for mainCollector
	mainCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Entering scraper::Scrape[mainCollector::OnRequest]")
		fmt.Printf("Visiting %s\n", r.URL)
		r.Headers.Set("content-type", "application/json; charset=utf-8")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Language", "en-ES,en;q=0.9,es-ES;q=0.8,es;q=0.7,en-GB;q=0.6,en-US;q=0.5")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Dnt", "1")
		r.Headers.Set("Priority", "u=0, i")
		r.Headers.Set("Sec-Cu-Ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
		r.Headers.Set("Set-Cu-Ua-Mobile", "?0")
		r.Headers.Set("Set-Cu-Ua-Platform", "MacOS")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	})
	mainCollector.OnHTML("li[class^='SearchPagestyled__Result-sc-'] a", func(e *colly.HTMLElement) {
		fmt.Println("Entering scraper::Scrape[mainCollector::OnHTML]")
		// This gets every URI to the first 20 properties matching the criteria.
		absolutePropertyUrl := e.Request.AbsoluteURL(e.Attr("href"))
		propertyCollector.Visit(absolutePropertyUrl)
		//propertyCollector.Wait()

		// We still need to navigate to "next page" and list every property.
	})
	mainCollector.OnScraped(func(r *colly.Response) {
		fmt.Println("Entering scraper::Scrape[mainCollector::OnScraped]")
	})
	mainCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Entering scraper::Scrape[mainCollector::OnError]")
		fmt.Printf("Error: %s", err)
	})

	// Define behaviour for propertyCollector
	propertyCollector.OnRequest(func(r *colly.Request) {
		fmt.Println("Entering scraper::Scrape[propertyCollector::OnRequest]")
		fmt.Printf("Visiting %s\n", r.URL)
		r.Headers.Set("content-type", "application/json; charset=utf-8")
		r.Headers.Set("Accept", "*/*")
		r.Headers.Set("Accept-Language", "en-ES,en;q=0.9,es-ES;q=0.8,es;q=0.7,en-GB;q=0.6,en-US;q=0.5")
		r.Headers.Set("Cache-Control", "max-age=0")
		r.Headers.Set("Dnt", "1")
		r.Headers.Set("Priority", "u=0, i")
		r.Headers.Set("Sec-Cu-Ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
		r.Headers.Set("Set-Cu-Ua-Mobile", "?0")
		r.Headers.Set("Set-Cu-Ua-Platform", "MacOS")
		r.Headers.Set("Sec-Fetch-Dest", "document")
		r.Headers.Set("Sec-Fetch-Mode", "navigate")
		r.Headers.Set("Sec-Fetch-Site", "none")
		r.Headers.Set("Sec-Fetch-User", "?1")
		r.Headers.Set("Upgrade-Insecure-Requests", "1")
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	})
	propertyCollector.OnHTML("div[class^='styles__MainColumn-sc-']", func(e *colly.HTMLElement) {
		fmt.Println("Entering scraper::Scrape[propertyCollector::OnHTML]")

		// Some links go to "developments" containing several "properties".
		// In that case, visit each property individually.
		subUnitsText := e.ChildText("div[data-testid='sub-units'] h3")
		if subUnitsText != "" {
			e.ForEach("div[data-testid='sub-units'] a[data-testid='sub-unit']", func(i int, e *colly.HTMLElement) {
				absolutePropertyUrl := e.Request.AbsoluteURL(e.Attr("href"))
				propertyCollector.Visit(absolutePropertyUrl)
			})
			return
		}

		property := model.Property{}

		property.Url = e.Request.URL.String()
		property.Address = e.ChildText("h1[data-testid='address']")
		property.Price, _ = extractPrice(e.ChildText("div[data-testid='price'] h2"))
		property.Type = e.ChildText("p[data-testid='property-type']")
		e.ForEach("div[data-testid='overview'] ul li", func(i int, e *colly.HTMLElement) {
			infoKey := e.ChildText("span")
			switch infoKey {
			case "Single Bedroom":
				property.NumSingleBedrooms, _ = extractNumSingleBedrooms(e.Text)
			case "Double Bedroom":
				property.NumDoubleBedrooms, _ = extractNumDoubleBedrooms(e.Text)
			case "Bathroom":
				property.NumBathrooms, _ = extractNumBathrooms(e.Text)
			case "Available From":
				property.AvailableFrom = extractAvailableFrom(e.Text)
			case "Furnished":
				property.Furnished = extractFurnished(e.Text)
			case "Lease":
				property.LeaseType = extractLeaseType(e.Text)
			}
		})
		property.Description = e.ChildText("div[data-testid='description'] div[data-testid='description']")
		property.ListingId = path.Base(property.Url)

		// Store property
		propertyIdUrlMap[property.ListingId] = property
	})
	propertyCollector.OnScraped(func(r *colly.Response) {
		fmt.Println("Entering scraper::Scrape[propertyCollector::OnScraped]")
	})
	propertyCollector.OnError(func(r *colly.Response, err error) {
		fmt.Println("Entering scraper::Scrape[propertyCollector::OnError]")
		fmt.Printf("Error: %s", err)
	})

	url, err := criteria.buildQuery(HTTP_SECURE + BASE_URL_DAFT)
	if err != nil {
		return nil
	}

	fmt.Println("Reconstructed URL from criteria: ", url)

	// Entry point to start collecting/scraping
	mainCollector.Visit(url)

	return propertyIdUrlMap
}

// Input format: "€X,XXX per month" / "€XXX per month" / "€XXX per week"
func extractPrice(textPrice string) (int, error) {
	sr := strings.NewReplacer("€", "", ",", "", " ", "", "per", "", "week", "", "month", "")
	textPrice = sr.Replace(textPrice)

	return utils.StringToInt(textPrice)
}

func extractNumSingleBedrooms(textNumSingleBedrooms string) (int, error) {
	textNumSingleBedrooms = strings.Replace(textNumSingleBedrooms, "Single Bedroom: ", "", 1)
	return utils.StringToInt(textNumSingleBedrooms)
}

func extractNumDoubleBedrooms(textNumDoubleBedrooms string) (int, error) {
	textNumDoubleBedrooms = strings.Replace(textNumDoubleBedrooms, "Double Bedroom: ", "", 1)
	return utils.StringToInt(textNumDoubleBedrooms)
}

func extractNumBathrooms(textNumBathrooms string) (int, error) {
	textNumBathrooms = strings.Replace(textNumBathrooms, "Bathroom: ", "", 1)
	return utils.StringToInt(textNumBathrooms)
}

func extractAvailableFrom(textAvailableFrom string) string {
	return strings.Replace(textAvailableFrom, "Available From: ", "", 1)
}

func extractFurnished(textFurnished string) bool {
	textFurnished = strings.Replace(textFurnished, "Furnished: ", "", 1)
	return textFurnished == "Yes"
}

func extractLeaseType(textLeaseType string) string {
	return strings.Replace(textLeaseType, "Lease: ", "", 1)
}
