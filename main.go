package main

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dev-drprasad/rest-api/auth"
	"github.com/dev-drprasad/rest-api/scraper"
	"github.com/dev-drprasad/rest-api/utils"
	"github.com/gorilla/mux"
	"gopkg.in/yaml.v2"
)

const salt = "salt"

type handler func(w http.ResponseWriter, r *http.Request)

func basicAuth(pass handler) handler {

	return func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

		if len(authHeader) != 2 || authHeader[0] != "Bearer" {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		h := sha1.New()
		h.Write([]byte(os.Getenv("HOODIE_PASSWORD")))
		sha1Password := hex.EncodeToString(h.Sum(nil))
		fmt.Println(sha1Password)
		hash := auth.GenerateBcryptHash(salt, os.Getenv("HOODIE_PASSWORD"))
		if hash != authHeader[1] {
			// If the two passwords don't match, return a 401 status
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
	}
	password := r.Form.Get("password")
	if password == os.Getenv("HOODIE_PASSWORD") {
		hash := auth.GenerateBcryptHash(salt, os.Getenv("HOODIE_PASSWORD"))
		utils.Respond(w, http.StatusOK, map[string]string{"token": hash}, nil)
	}
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

func loadDefinitionFromFile(filename string) *scraper.Definition {
	f, err := ioutil.ReadFile(path.Join("definitions", filename))
	if err != nil {
		log.Printf("Failed to read definition from file=%s error=%s", filename, err)
		return nil
	}
	d := &scraper.Definition{}
	err = yaml.Unmarshal(f, &d)
	if err != nil {
		log.Fatalf("Failed to Unmarshal file=%s error=%v", filename, err)
		return nil
	}
	d.ID = strings.TrimSuffix(filename, path.Ext(filename))

	return d
}

func getSiteDefinition(site string) *scraper.Definition {
	return loadDefinitionFromFile(site + ".yaml")
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
		utils.Respond(w, http.StatusBadRequest, nil, errors.New("Missing 'site' in queryparams"))
		return
	}
	if query == "" {
		utils.Respond(w, http.StatusBadRequest, nil, errors.New("Missing 'query' in queryparams"))
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

func handleSourceList(w http.ResponseWriter, r *http.Request) {
	defRoot := "definitions"
	files, err := ioutil.ReadDir(defRoot)
	if err != nil {
		log.Printf("Failed to read definitions directory Error=%s", err.Error())
		utils.Respond(w, http.StatusInternalServerError, nil, err)
		return
	}

	sources := []map[string]string{}
	for _, f := range files {
		if !f.IsDir() {
			def := loadDefinitionFromFile(f.Name())
			if def == nil {
				log.Printf("Failed to get definition for source=%s", f.Name())
				continue
			}

			sources = append(sources, map[string]string{"id": def.ID, "name": def.Name})
		}
	}

	utils.Respond(w, http.StatusOK, sources, nil)
}

func main() {
	err := scraper.Init()
	if err != nil {
		log.Println("Error in initializing scraper")
		log.Println(err)
	}

	router := mux.NewRouter()

	// router.HandleFunc("/mag2tor", handleMag2Tor).Methods("GET")
	router.HandleFunc("/api/v1/search", basicAuth(search)).Methods("GET")
	router.HandleFunc("/api/v1/detail", basicAuth(handleDetailScrape)).Methods("GET")
	router.HandleFunc("/api/v1/list", basicAuth(handleListScrape)).Methods("GET")
	router.HandleFunc("/api/v1/sources", basicAuth(handleSourceList)).Methods("GET")
	router.HandleFunc("/api/v1/login", login).Methods("POST")

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./client/")))
	// router.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	http.ServeFile(w, r, "./index.html")
	// })

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
