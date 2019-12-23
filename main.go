package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/dev-drprasad/rest-api/scraper"
	"github.com/dev-drprasad/rest-api/utils"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

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

func getSiteDefinition(site string) *scraper.Definition {
	f, err := ioutil.ReadFile("definitions/" + site + ".yaml")
	if err != nil {
		log.Printf("error reading definition %s: %s", site, err)
		return nil
	}
	d := &scraper.Definition{}
	err = yaml.Unmarshal(f, &d)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
		return nil
	}
	d.ID = site
	log.Printf("%v", d)
	return d
}

func handleListScrape(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	listURL := r.URL.Query().Get("url")

	if listURL == "" {
		utils.Respond(w, http.StatusBadRequest, nil, errors.New("Missing 'url' in queryparams"))
		return
	}
	if site == "" {
		utils.Respond(w, http.StatusBadRequest, nil, errors.New("Missing 'site' in queryparams"))
		return
	}

	// TODO : Check `site` is valid
	config := getSiteDefinition(site)
	ti, err := scraper.ScrapeList(*config, listURL)

	if err != nil {
		utils.Respond(w, http.StatusInternalServerError, nil, err)
		return
	}

	utils.Respond(w, http.StatusOK, ti, nil)
	return
}

func handleDetailScrape(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	detailURL := r.URL.Query().Get("url")

	if detailURL == "" {
		utils.Respond(w, http.StatusBadRequest, nil, errors.New("No value present for queryparam 'url'"))
		return
	}
	if site == "" {
		utils.Respond(w, http.StatusBadRequest, nil, errors.New("No value present for queryparam 'site'"))
		return
	}

	config := getSiteDefinition(site)
	ti, err := scraper.ScrapeDetail(*config, detailURL)

	if err != nil {
		utils.Respond(w, http.StatusInternalServerError, nil, err)
		return
	}

	utils.Respond(w, http.StatusOK, ti, nil)
	return
}

func search(w http.ResponseWriter, r *http.Request) {
	site := r.URL.Query().Get("site")
	query := r.URL.Query().Get("query")
	pageNo := r.URL.Query().Get("pageNo")

	if site == "" {
		utils.Respond(w, http.StatusInternalServerError, nil, errors.New("Missing 'site' in queryparams"))
		return
	}
	if query == "" {
		utils.Respond(w, http.StatusInternalServerError, nil, errors.New("Missing 'query' in queryparams"))
		return
	}

	d := getSiteDefinition(site)

	var relativeURL string
	if pageNo != "" {
		relativeURL = utils.ProcessString(d.Search.Pagination.URL, map[string]string{"page": pageNo, "query": query})
	} else {
		relativeURL = utils.ProcessString(d.Search.Path, map[string]string{"query": query})
	}

	if relativeURL == "" {
		err := errors.New("Relative URL for list is empty")
		log.Printf(err.Error())
		utils.Respond(w, http.StatusInternalServerError, nil, err)
	}

	url, _ := url.Parse(d.URL)
	url.Path = path.Join(url.Path, relativeURL) + "/" //1337x need '/' at end

	til, err := scraper.ScrapeList(*d, url.String())
	if err != nil {
		utils.Respond(w, http.StatusInternalServerError, nil, err)
		return
	}

	utils.Respond(w, http.StatusOK, til, nil)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)

	// router.HandleFunc("/mag2tor", handleMag2Tor).Methods("GET")
	router.HandleFunc("/api/v1/search", search).Methods("GET")
	router.HandleFunc("/api/v1/detail", handleDetailScrape).Methods("GET")
	router.HandleFunc("/api/v1/list", handleListScrape).Methods("GET")

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      router, // Pass our instance of gorilla/mux in.
	}

	log.Println("Starting app...")
	log.Fatal(srv.ListenAndServe())
}
