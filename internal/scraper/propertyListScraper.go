package scraper

import (
	"fmt"
	"strings"

	"github.com/gocolly/colly"
)

// const HTTP_SECURE string = "https://"
// const BASE_URL_DAFT string = "www.daft.ie"

// type Scraper interface {
// 	Scrape()
// }

type PropertyListScraper struct {
	//propertyIds []string
}

func (self *PropertyListScraper) Scrape(url string) {

	c_list := colly.NewCollector(
		colly.AllowedDomains( /*"daft.ie",*/ BASE_URL_DAFT),
		colly.DetectCharset(),
	)

	c_list.OnRequest(func(r *colly.Request) {
		//r.Headers.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
		r.Headers.Set("content-type", "application/json; charset=utf-8")
		r.Headers.Set("Accept", "*/*")
		//r.Headers.Set("Accept-Encoding", "gzip, deflate, br, zstd")
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
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c_list.OnHTML("li[class^='SearchPagestyled__Result-sc-'] a", func(h *colly.HTMLElement) {
		// This gets every URI to the first 20 properties matching the criteria.
		fmt.Println("OnHTML")
		property_url := BASE_URL_DAFT + h.Attr("href")
		fmt.Println(property_url)

		// TODO follow here!
		// TODO use a formatter here
		parts := strings.Split(property_url, "/")
		lastPart := parts[len(parts)-1:]
		fmt.Println(lastPart) // []string
		res := fmt.Sprint(lastPart)
		fmt.Println(res) // string

		// c_property.Visit(property_url)

		// We still need to navigate to "next page" and list every property.
	})

	c_list.OnScraped(func(r *colly.Response) {

	})

	c_list.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error: %s", err)
	})

	//c_list.Visit(HTTP_SECURE + BASE_URL_DAFT + "/property-for-rent/dublin-9-dublin?rentalPrice_to=2000")
	c_list.Visit(HTTP_SECURE + BASE_URL_DAFT + url)

}
