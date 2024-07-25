package main

import (
	"fmt"

	"github.com/venturarome/DaftWatch/internal/server"
)

func main() {

	// cc := colly.NewCollector(
	// 	colly.AllowedDomains( /*"daft.ie",*/ "www.daft.ie"),
	// 	colly.DetectCharset(),
	// )
	// cc.OnRequest(func(r *colly.Request) {
	// 	r.Headers.Set("content-type", "application/json; charset=utf-8")
	// 	r.Headers.Set("Accept", "*/*")
	// 	r.Headers.Set("Accept-Language", "en-ES,en;q=0.9,es-ES;q=0.8,es;q=0.7,en-GB;q=0.6,en-US;q=0.5")
	// 	r.Headers.Set("Cache-Control", "max-age=0")
	// 	r.Headers.Set("Dnt", "1")
	// 	r.Headers.Set("Priority", "u=0, i")
	// 	r.Headers.Set("Sec-Cu-Ua", `"Not/A)Brand";v="8", "Chromium";v="126", "Google Chrome";v="126"`)
	// 	r.Headers.Set("Set-Cu-Ua-Mobile", "?0")
	// 	r.Headers.Set("Set-Cu-Ua-Platform", "MacOS")
	// 	r.Headers.Set("Sec-Fetch-Dest", "document")
	// 	r.Headers.Set("Sec-Fetch-Mode", "navigate")
	// 	r.Headers.Set("Sec-Fetch-Site", "none")
	// 	r.Headers.Set("Sec-Fetch-User", "?1")
	// 	r.Headers.Set("Upgrade-Insecure-Requests", "1")
	// 	r.Headers.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	// 	fmt.Printf("Visiting %s\n", r.URL)
	// })
	// cc.OnScraped(func(r *colly.Response) {
	// 	fmt.Println("Scraped main page!")
	// })
	// cc.OnError(func(r *colly.Response, err error) {
	// 	fmt.Printf("Error: %s", err)
	// })
	// cc.OnHTML("div[class^='styles__MainColumn-sc-']", func(e *colly.HTMLElement) {
	// 	property := scraper.Property{}
	// 	// ADDRESS
	// 	address := e.ChildText("h1[data-testid='address']")
	// 	fmt.Println("Address: ", address)
	// 	property.address = address

	// 	// PRICE
	// 	priceText := e.ChildText("div[data-testid='price'] h2") // "€1,800 per month"
	// 	price, err := strconv.ParseUint(priceText, 10, 64)
	// 	if err != nil {
	// 		price = 0
	// 	}
	// 	//property.price = uint16(price)
	// 	fmt.Println("Price: ", price)

	// 	// OVERVIEW
	// 	e.ForEach("div[data-testid='overview'] ul li", func(i int, e *colly.HTMLElement) {
	// 		fmt.Println(e.ChildText("span"))
	// 		fmt.Println("Overview info: ", e.Text)

	// 		// item := Product{}
	// 		// item.Name = h.Text
	// 		// item.Image = e.ChildAttr("img", "data-src")
	// 		// item.Price = e.Attr("data-price")
	// 		// item.Url = "https://jumia.com.ng" + e.Attr("href")
	// 		// item.Discount = e.ChildText("div.tag._dsct")
	// 		// products = append(products, item)
	// 	})

	// 	// DESCRIPTION
	// 	fmt.Println("Description: ", e.ChildText("div[data-testid='description'] div[data-testid='description']"))

	// })
	// cc.Visit("https://www.daft.ie/for-rent/apartment-9-iona-road-glasnevin-dublin-9/5701754")

	server := server.NewServer()

	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
