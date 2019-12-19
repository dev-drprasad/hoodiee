package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"time"

	"github.com/dev-drprasad/rest-api/utils"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

const userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.3998.0 Safari/537.36"
const scrapeDetail = false

type Selector struct {
	Selector string `yaml:"selector"`
	Attr     string `yaml:"attr"`
}

type Fields struct {
	Name      Selector `yaml:"name"`
	MagnetURI Selector `yaml:"magnet"`
	URL       Selector `yaml:"url"`
	Seeders   Selector `yaml:"seeders"`
	Leechers  Selector `yaml:"leechers"`
}

type Definition struct {
	Name   string `yaml:"name"`
	URL    string `yaml:"url"`
	Search struct {
		Path string `yaml:"path"`
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

type TorrentInfo struct {
	Name      string `json:"name"`
	MagnetURI string `json:"magnetURI"`
	URL       string `json:"URL"`
	Seeders   string `json:"seeders"`
	Leechers  string `json:"leechers"`
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

// func handleMag2Tor(w http.ResponseWriter, r *http.Request) {
// 	magnetURI := r.URL.Query().Get("magnetURI")
// 	log.Printf("magnetURI: %s", magnetURI)
// 	filename := mag2tor.Mag2Tor(magnetURI)
// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		log.Printf("error reading file %s, %s", filename, err)
// 	}
// 	w.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
// 	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data))
// }

// func tpb(w http.ResponseWriter, r *http.Request) {
// 	c := colly.NewCollector()

// 	// Find and visit all links
// 	c.OnHTML("table#searchResult", func(e *colly.HTMLElement) {
// 		m := []TorrentInfo{}
// 		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
// 			magnetURL := e.ChildAttr("td:nth-child(2) a[href^=\"magnet:\"]", "href")
// 			name := e.ChildText("td:nth-child(2) div.detName")
// 			log.Println(name, magnetURL)
// 			m = append(m, TorrentInfo{Name: name, MagnetURI: magnetURL})
// 		})
// 		json.NewEncoder(w).Encode(m)
// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		log.Println("Visiting", r.URL)
// 		r.Headers.Set("Accept", "text/html")
// 		r.Headers.Set("Accept-Encoding", "gzip")
// 	})

// 	log.Println("visiting URL...")
// 	c.Visit("https://pirateproxy.onl/search/tomb%20raider/")
// }

// func kickasstorrent(w http.ResponseWriter, r *http.Request) {
// 	c := colly.NewCollector()
// 	// c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.3998.0 Safari/537.36"
// 	c.UserAgent = "Mozilla/5.0 (Linux; Android 9; SM-G960F Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.157 Mobile Safari/537.36"
// 	dc := c.Clone()

// 	dc.Limit(&colly.LimitRule{
// 		Parallelism: 2,
// 		Delay:       2 * time.Second,
// 	})

// 	m := []TorrentInfo{}
// 	// Find and visit all links
// 	dc.OnHTML("table#mainDetailsTable", func(e *colly.HTMLElement) {
// 		magnetURL := e.ChildAttr(".downloadButtonGroup a[href^=\"magnet:\"]", "href")
// 		name := e.ChildText("h1")
// 		log.Println(name, magnetURL)
// 		m = append(m, TorrentInfo{Name: name, MagnetURI: magnetURL})
// 	})

// 	c.OnHTML("table.frontPageWidget", func(e *colly.HTMLElement) {
// 		// log.Println(e.Text)
// 		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
// 			// log.Println(string(e.Response.Body))
// 			rURL := e.ChildAttr("td:nth-child(1) div.torrentname a.cellMainLink", "href")
// 			detailURL := e.Request.AbsoluteURL(rURL)
// 			dc.Visit(detailURL)
// 		})

// 	})

// 	c.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 		r.Headers.Set("Accept", "text/html")
// 		r.Headers.Set("Accept-Encoding", "gzip")

// 	})

// 	c.OnError(func(r *colly.Response, err error) {
// 		log.Printf("response: %s", string(r.Body))
// 		log.Println("error:", r.StatusCode, err)
// 		w.WriteHeader(500)
// 	})

// 	c.OnScraped(func(r *colly.Response) {
// 		log.Println("Completed")
// 		json.NewEncoder(w).Encode(m)
// 	})

// 	dc.OnError(func(r *colly.Response, err error) {
// 		log.Println("error:", r.StatusCode, err)
// 	})

// 	dc.OnRequest(func(r *colly.Request) {
// 		fmt.Println("Visiting", r.URL)
// 	})

// 	log.Println("visiting URL...")
// 	c.Visit("https://katcr.to/usearch/game%20of")
// }

// func (f )

func getSiteDefinition(site string) *Definition {
	f, err := ioutil.ReadFile("definitions/" + site + ".yaml")
	if err != nil {
		log.Printf("error reading definition %s: %s", site, err)
		return nil
	}
	d := &Definition{}
	err = yaml.Unmarshal(f, &d)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil
	}
	log.Printf("%v", d)
	return d
}

func scrapeDetails(site string, detailURL string) (TorrentInfo, error) {
	config := getSiteDefinition(site)
	ti := TorrentInfo{}

	if config.Search.Detail == nil || config.Search.Detail.Fields == nil {
		return ti, errors.New("No Fields defined in config file")
	}

	scraper := colly.NewCollector()
	scraper.UserAgent = userAgent

	var err error
	scraper.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	scraper.OnError(func(r *colly.Response, e error) {
		log.Println("error:", r.StatusCode, e)
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
						log.Println("setting...")
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
	log.Println("returning...")

	return ti, err
}

func handleDetailScrape(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	detailURL := r.URL.Query().Get("url")

	ti, err := scrapeDetails(site, detailURL)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	Respond(w, ti, nil)
}

func search(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	query := r.URL.Query().Get("query")
	if site == "" || query == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	d := getSiteDefinition(site)

	u, _ := url.Parse(d.URL)
	searchURL := utils.ProcessString(d.Search.Path, map[string]string{"query": query})
	searchAbsURL := d.URL + searchURL
	log.Printf("searchAbsURL %s", searchAbsURL)

	m := []TorrentInfo{}

	c := colly.NewCollector()
	// c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.3998.0 Safari/537.36"
	// c.UserAgent = "Mozilla/5.0 (Linux; Android 9; SM-G960F Build/PPR1.180610.011; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/74.0.3729.157 Mobile Safari/537.36"
	c.UserAgent = userAgent

	log.Printf("host: %s", u.Host)

	var dc *colly.Collector
	if scrapeDetail && d.Search.Detail != nil {
		dc = c.Clone()

		dc.Limit(&colly.LimitRule{
			Parallelism: 2,
			Delay:       3 * time.Second,
		})

		dc.OnRequest(func(r *colly.Request) {
			fmt.Println("Visiting", r.URL)
		})

		dc.OnError(func(r *colly.Response, err error) {
			log.Println("error:", r.StatusCode, err)
			w.WriteHeader(500)
		})

	}

	c.OnHTML("body", func(e *colly.HTMLElement) {
		e.ForEach(d.Search.List.Selector, func(i int, e *colly.HTMLElement) {

			fieldValues := reflect.ValueOf(d.Search.List.Fields)
			typeOfDefinition := fieldValues.Type()

			ti := TorrentInfo{}
			tis := reflect.ValueOf(&ti).Elem()

			for i := 0; i < fieldValues.NumField(); i++ {
				sel := fieldValues.Field(i).Interface().(Selector)
				fieldName := typeOfDefinition.Field(i).Name
				log.Printf("%v", sel)
				log.Printf("%s", fieldName)

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
			log.Printf("%#v", ti)

			m = append(m, ti)
			idx := len(m) - 1

			if scrapeDetail && d.Search.Detail != nil {
				dc.OnHTML("body", func(e *colly.HTMLElement) {

					if d.Search.Detail.Fields.MagnetURI.Selector != "" {
						var magnetURI string
						if d.Search.Detail.Fields.MagnetURI.Attr != "" {
							magnetURI = e.ChildAttr(d.Search.Detail.Fields.MagnetURI.Selector, d.Search.Detail.Fields.MagnetURI.Attr)
						} else {
							magnetURI = e.ChildText(d.Search.Detail.Fields.MagnetURI.Selector)
						}

						m[idx].MagnetURI = magnetURI
					}

					if d.Search.Detail.Fields.Name.Selector != "" {
						var name string
						if d.Search.Detail.Fields.Name.Attr != "" {
							name = e.ChildAttr(d.Search.Detail.Fields.Name.Selector, d.Search.Detail.Fields.Name.Attr)
						} else {
							name = e.ChildText(d.Search.Detail.Fields.Name.Selector)
						}
						m[idx].Name = name
					}
					log.Printf("%d", idx)
					// log.Printf("%v", m[idx])
				})
				log.Printf("visiting detail... %d", idx)
				dc.Visit(ti.URL)
				log.Printf("detaching detail %d", idx)
				dc.OnHTMLDetach("body")
			}
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
		r.Headers.Set("Accept", "text/html")
		r.Headers.Set("Accept-Encoding", "gzip")
		r.Headers.Set("Host", u.Host)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request Error: ", err, "Code: ", r.StatusCode)
		w.WriteHeader(500)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Completed")
		// json.NewEncoder(w).Encode(m)
		Respond(w, m, nil)
	})

	log.Printf("%v", d.Search.Detail)

	c.Visit(searchAbsURL)
}

func Respond(w http.ResponseWriter, data interface{}, err error) {
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data, "error": err})
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/tpb", tpb).Methods("GET")
	// router.HandleFunc("/kat", kickasstorrent).Methods("GET")
	// router.HandleFunc("/mag2tor", handleMag2Tor).Methods("GET")
	router.HandleFunc("/api/v1/search", search).Methods("GET")
	router.HandleFunc("/api/v1/detail", handleDetailScrape).Methods("GET")
	log.Println("Starting app...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
