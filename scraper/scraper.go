package scraper

import (
	"errors"
	"fmt"
	"log"
	"net/http/cookiejar"
	"net/url"
	"reflect"
	"regexp"
	"strconv"

	"github.com/gocolly/colly"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.3998.0 Safari/537.36"

// const userAgent = "Mozilla/5.0 (Linux; Android 9; SM-G960F Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.157 Mobile Safari/537.36"
const scrapeDetail = false

type TorrentInfo struct {
	Name      string `json:"name"`
	MagnetURI string `json:"magnetURI"`
	URL       string `json:"URL"`
	Seeders   string `json:"seeders"`
	Leechers  string `json:"leechers"`
	Source    string `json:"source"`
}

type Transformer struct {
	Type  string `yaml:"type"`
	Value string `yaml:"value"`
}

type ScrapeResult struct {
	Source string        `json:"source"`
	Pages  int           `json:"pages"`
	Items  []TorrentInfo `json:"items"`
}

type Selector struct {
	Selector  string        `yaml:"selector"`
	Attr      string        `yaml:"attr"`
	Transform []Transformer `yaml:"transform"`
}

type Fields struct {
	Name      Selector `yaml:"name"`
	MagnetURI Selector `yaml:"magnet"`
	URL       Selector `yaml:"url"`
	Seeders   Selector `yaml:"seeders"`
	Leechers  Selector `yaml:"leechers"`
}

type Definition struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	Search struct {
		Path       string `yaml:"path"`
		Pagination struct {
			URL   string   `yaml:"url"`
			Total Selector `yaml:"total"`
		} `yaml:"pagination"`
		List struct {
			Fields   Fields `yaml:"fields"`
			Selector string `yaml:"selector"`
			Attr     string `yaml:"attr"`
		} `yaml:"list"`
		Detail *struct {
			Fields *Fields `yaml:"fields"`
		} `yaml:"details"`
	} `yaml:"search"`
}

var jar *cookiejar.Jar

func Init() error {
	var err error
	if jar == nil {
		jar, err = cookiejar.New(nil)
	}
	return err
}

func NewScrapeResult() ScrapeResult {
	return ScrapeResult{
		Pages: 1,
		Items: []TorrentInfo{},
	}
}

func NewScraper() *colly.Collector {
	scraper := colly.NewCollector()
	scraper.SetCookieJar(jar)
	scraper.UserAgent = userAgent
	return scraper
}

func ScrapeList(config Definition, listURL string) (ScrapeResult, error) {
	u, _ := url.Parse(config.URL)

	var err error
	result := NewScrapeResult()

	scraper := NewScraper()

	scraper.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach(config.Search.List.Selector, func(i int, e *colly.HTMLElement) {

			fieldValues := reflect.ValueOf(config.Search.List.Fields)
			typeOfDefinition := fieldValues.Type()

			ti := TorrentInfo{Source: config.ID}
			tis := reflect.ValueOf(&ti).Elem()

			for i := 0; i < fieldValues.NumField(); i++ {
				sel := fieldValues.Field(i).Interface().(Selector)
				fieldName := typeOfDefinition.Field(i).Name

				var parsed string
				if sel.Attr != "" {
					parsed = e.ChildAttr(sel.Selector, sel.Attr)
				} else {
					parsed = e.ChildText(sel.Selector)
				}

				// exported field
				tif := tis.FieldByName(fieldName)
				if tif.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					if tif.CanSet() {
						// change value of N
						switch tif.Kind() {
						case reflect.Int:
							pn, _ := strconv.Atoi(parsed)
							tif.SetInt(int64(pn))
						case reflect.String:
							tif.SetString(parsed)
						default:
						}
					}
				} else {
					log.Printf("Field '%s' is not valid", fieldName)
				}

			}

			ti.URL = e.Request.AbsoluteURL(ti.URL)

			result.Items = append(result.Items, ti)
		})

		log.Printf("Scraped %d items from list page of %s", len(result.Items), config.Name)

		if config.Search.Pagination.Total.Selector != "" {
			var pagesStr string
			if config.Search.Pagination.Total.Attr != "" {
				pagesStr = e.ChildAttr(config.Search.Pagination.Total.Selector, config.Search.Pagination.Total.Attr)
			} else {
				pagesStr = e.ChildText(config.Search.Pagination.Total.Selector)
			}
			log.Printf("Total num of pages: %s", pagesStr)

			if len(config.Search.Pagination.Total.Transform) > 0 {
				transforms := config.Search.Pagination.Total.Transform

				for _, transform := range transforms {
					if transform.Type == "regex" {
						regx := regexp.MustCompile(transform.Value)
						match := regx.FindStringSubmatch(pagesStr)
						if len(match) == 0 {
							log.Printf("No match found for pages=%s with given regex=%s", pagesStr, transform.Value)
							return
						}
						if len(match) == 1 {
							log.Printf("No submatch found for pages=%s with given regex=%s", pagesStr, transform.Value)
							return
						}
						if len(match) > 1 {
							pagesStr = match[1]
						}
					}
				}
			}

			pages, err := strconv.Atoi(pagesStr)
			if err != nil {
				log.Printf("Parsing total num of pages failed. Received=%s", pagesStr)
				return
			}
			result.Pages = pages
		} else {
			log.Printf("No selector found for total num of pages")
		}
	})

	scraper.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", r.URL)
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Accept-Encoding", "gzip")
		r.Headers.Set("Host", u.Host)
		log.Printf("Cookie: %s", r.Headers.Get("Cookie"))
	})

	scraper.OnError(func(r *colly.Response, e error) {
		log.Printf("Request Error: Error=%s Code=%d", e.Error(), r.StatusCode)
		err = e
	})

	scraper.Visit(listURL)
	result.Source = config.Name
	return result, err
}

func ScrapeDetail(config Definition, detailURL string) (*TorrentInfo, error) {

	ti := TorrentInfo{Source: config.ID}

	if config.Search.Detail.Fields == nil {
		return nil, errors.New("No Fields defined in config file")
	}

	scraper := NewScraper()

	var err error
	scraper.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	scraper.OnError(func(r *colly.Response, e error) {
		log.Printf("Request Error: Error=%s Code=%d", e.Error(), r.StatusCode)
		err = e
	})

	scraper.OnHTML("body", func(e *colly.HTMLElement) {

		fieldValues := reflect.Indirect(reflect.ValueOf(config.Search.Detail.Fields))
		typeOfDefinition := fieldValues.Type()

		for i := 0; i < fieldValues.NumField(); i++ {
			sel := fieldValues.Field(i).Interface().(Selector)
			fieldName := typeOfDefinition.Field(i).Name
			log.Println("field: ", fieldName, "selector: ", sel)

			var parsed string
			if sel.Selector != "" {

				if sel.Attr != "" {
					parsed = e.ChildAttr(sel.Selector, sel.Attr)
				} else {
					parsed = e.ChildText(sel.Selector)
				}
			} else {
				log.Println("selector is empty for ", fieldName)
			}
			log.Println("parsed: ", parsed)

			tis := reflect.ValueOf(&ti).Elem()
			// exported field
			tif := tis.FieldByName(fieldName)
			if tif.IsValid() {
				// A Value can be changed only if it is
				// addressable and was not obtained by
				// the use of unexported struct fields.
				if tif.CanSet() {
					// change value of N
					switch tif.Kind() {
					case reflect.Int:
						pn, _ := strconv.Atoi(parsed)
						tif.SetInt(int64(pn))
					case reflect.String:
						tif.SetString(parsed)
					default:
					}
				}
			} else {
				log.Printf("Field '%s' is not valid", fieldName)
			}

		}
	})

	scraper.Visit(detailURL)

	return &ti, err
}
