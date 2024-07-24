package scraper

import (
	"fmt"

	"github.com/gocolly/colly"
)

type PropertyScraper struct {
}

func (self *PropertyScraper) Scrape(url string) {

	// propertiesInfo := []Property

	c_property := colly.NewCollector(
		colly.AllowedDomains( /*"daft.ie",*/ BASE_URL_DAFT),
		colly.DetectCharset(),
	)
	c_property.OnRequest(func(r *colly.Request) {
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
	c_property.OnHTML("h1[data-testid='address']", func(h *colly.HTMLElement) {
		// This gets ADDRESS information from a property.
		fmt.Println("c_property::OnHTML::address")
		fmt.Println(h.Text)
	})
	c_property.OnHTML("div[data-testid='price'] h2", func(h *colly.HTMLElement) {
		// This gets PRICE information from a property.
		fmt.Println("c_property::OnHTML::price")
		fmt.Println(h.Text)
	})
	c_property.OnHTML("div[data-testid='overview']", func(h *colly.HTMLElement) {
		// This gets OVERVIEW information from a property (bedrooms, bathrooms, availability, furnished, lease, ...).
		fmt.Println("c_property::OnHTML::overview")
		h.ForEach("ul li", func(i int, h *colly.HTMLElement) {
			fmt.Println(h.Text)
		})
	})
	c_property.OnHTML("div[data-testid='description'] div[data-testid='description']", func(h *colly.HTMLElement) {
		// This gets DESCRIPTION information from a property.
		fmt.Println("c_property::OnHTML::description")
		fmt.Println(h.Text)
	})
	c_property.OnScraped(func(r *colly.Response) {

	})
	c_property.OnError(func(r *colly.Response, err error) {
		fmt.Printf("Error: %s", err)
	})
	c_property.Visit(HTTP_SECURE + BASE_URL_DAFT + "/for-rent/apartment-9-iona-road-glasnevin-dublin-9/5701754")
}
