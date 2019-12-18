package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dev-drprasad/rest-api/mag2tor"
	"github.com/gocolly/colly"
	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

type TorrentResult struct {
	Name      string `json:"name"`
	MagnetURL string `json:"magnetURL"`
}

type allEvents []event

var events = allEvents{
	{
		ID:          "1",
		Title:       "Introduction to Golang",
		Description: "Come join us for a chance to learn how golang works and get to eventually try it out",
	},
}

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome home!")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}

	json.Unmarshal(reqBody, &newEvent)
	events = append(events, newEvent)
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(newEvent)
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for _, singleEvent := range events {
		if singleEvent.ID == eventID {
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(events)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]
	var updatedEvent event

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(w, "Kindly enter data with the event title and description only in order to update")
	}
	json.Unmarshal(reqBody, &updatedEvent)

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			singleEvent.Title = updatedEvent.Title
			singleEvent.Description = updatedEvent.Description
			events = append(events[:i], singleEvent)
			json.NewEncoder(w).Encode(singleEvent)
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	eventID := mux.Vars(r)["id"]

	for i, singleEvent := range events {
		if singleEvent.ID == eventID {
			events = append(events[:i], events[i+1:]...)
			fmt.Fprintf(w, "The event with ID %v has been deleted successfully", eventID)
		}
	}
}

func handleMag2Tor(w http.ResponseWriter, r *http.Request) {
	magnetURI := r.URL.Query().Get("magnetURI")
	log.Printf("magnetURI: %s", magnetURI)
	filename := mag2tor.Mag2Tor(magnetURI)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("error reading file %s, %s", filename, err)
	}
	w.Header().Add("Content-Disposition", "attachment; filename=\""+filename+"\"")
	http.ServeContent(w, r, filename, time.Now(), bytes.NewReader(data))
}

func tpb(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector()

	// Find and visit all links
	c.OnHTML("table#searchResult", func(e *colly.HTMLElement) {
		m := []TorrentResult{}
		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
			magnetURL := e.ChildAttr("td:nth-child(2) a[href^=\"magnet:\"]", "href")
			name := e.ChildText("td:nth-child(2) div.detName")
			log.Println(name, magnetURL)
			m = append(m, TorrentResult{Name: name, MagnetURL: magnetURL})
		})
		json.NewEncoder(w).Encode(m)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	log.Println("visiting URL...")
	c.Visit("https://pirateproxy.onl/search/tomb%20raider/")
}

func kickasstorrent(w http.ResponseWriter, r *http.Request) {
	c := colly.NewCollector()
	c.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.3998.0 Safari/537.36"

	dc := c.Clone()

	dc.Limit(&colly.LimitRule{
		Parallelism: 2,
		Delay:       2 * time.Second,
	})

	m := []TorrentResult{}
	// Find and visit all links
	dc.OnHTML("table#mainDetailsTable", func(e *colly.HTMLElement) {
		magnetURL := e.ChildAttr(".downloadButtonGroup a[href^=\"magnet:\"]", "href")
		name := e.ChildText("h1")
		log.Println(name, magnetURL)
		m = append(m, TorrentResult{Name: name, MagnetURL: magnetURL})
	})

	c.OnHTML("table.frontPageWidget", func(e *colly.HTMLElement) {
		// log.Println(e.Text)
		e.ForEach("tr", func(i int, e *colly.HTMLElement) {
			// log.Println(string(e.Response.Body))
			rURL := e.ChildAttr("td:nth-child(1) div.torrentname a.cellMainLink", "href")
			detailURL := e.Request.AbsoluteURL(rURL)
			dc.Visit(detailURL)
		})

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		w.WriteHeader(500)
	})

	dc.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
	})

	dc.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnScraped(func(r *colly.Response) {
		log.Println("Completed")
		json.NewEncoder(w).Encode(m)
	})

	log.Println("visiting URL...")
	c.Visit("https://katcr.to/usearch/game%20of")
}

func main() {
	// initEvents()
	router := mux.NewRouter().StrictSlash(true)
	magnetLink := "magnet:?xt=urn:btih:bac2c9d9c552ab2465485fd37c11877f9af051db&dn=Rick and Morty S04E03 One Crew Over The Crewcoos Morty 1080p AMZN WEBRip DDP5 1 x264 CtrlHD [rartv]&tr=udp://tracker.coppersurfer.tk:6969&tr=udp://tracker.opentrackr.org:1337&tr=udp://tracker.pirateparty.gr:6969&tr=udp://9.rarbg.to:2710&tr=udp://9.rarbg.me:2710"
	// re := regexp.MustCompile(`xt=urn:btih:([^&/]+)`)
	// fmt.Printf("%q\n", re.Find([]byte(magnetLink)))
	tc := "d10:magnet-uri" + strconv.Itoa(len(magnetLink)) + ":" + magnetLink + "e"

	ioutil.WriteFile("filename.torrent", []byte(tc), 0644)

	router.HandleFunc("/", homeLink)
	router.HandleFunc("/event", createEvent).Methods("POST")
	router.HandleFunc("/events", getAllEvents).Methods("GET")
	router.HandleFunc("/events/{id}", getOneEvent).Methods("GET")
	router.HandleFunc("/events/{id}", updateEvent).Methods("PATCH")
	router.HandleFunc("/events/{id}", deleteEvent).Methods("DELETE")
	router.HandleFunc("/tpb", tpb).Methods("GET")
	router.HandleFunc("/kat", kickasstorrent).Methods("GET")
	router.HandleFunc("/mag2tor", handleMag2Tor).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", router))
}
