package scraper

const HTTP_SECURE string = "https://"
const BASE_URL_DAFT string = "www.daft.ie"

type Scraper interface {
	Scrape()
}
