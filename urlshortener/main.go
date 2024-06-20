//go:build !solution

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var (
	idCounter int
	idMutex   sync.Mutex
)

var (
	urlsToTiny map[string]string
	tinyToUrls map[string]string
	dbMutex    sync.Mutex
)

func main() {

	urlsToTiny = make(map[string]string)
	tinyToUrls = make(map[string]string)

	portPtr := flag.Int("port", -1, "port")

	flag.Parse()

	if *portPtr == -1 {
		panic("Missing port")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Post("/shorten", shortenHandler)
	r.Get("/go/{key}", getKeyHandler)

	log.Fatal(http.ListenAndServe(fmt.Sprintf("localhost:%d", *portPtr), r))
}

func shortenHandler(w http.ResponseWriter, r *http.Request) {
	var s ShortenRequestBody

	err := json.NewDecoder(r.Body).Decode(&s)
	if err != nil {
		w.WriteHeader(400)
		return
	}

	var tinyURL string

	dbMutex.Lock()
	value, ok := urlsToTiny[s.URL]
	dbMutex.Unlock()

	if ok {
		tinyURL = value
	} else {
		var id int

		idMutex.Lock()
		id = idCounter
		idCounter++
		idMutex.Unlock()

		tinyURL = strconv.Itoa(id)

		dbMutex.Lock()
		urlsToTiny[s.URL] = tinyURL
		tinyToUrls[tinyURL] = s.URL
		dbMutex.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, `{"url":"%s","key":"%s"}`, s.URL, tinyURL)
}

func getKeyHandler(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")

	dbMutex.Lock()
	value, ok := tinyToUrls[key]
	dbMutex.Unlock()

	if !ok {
		w.WriteHeader(404)
		fmt.Fprint(w, "key not found")
		return
	}

	w.Header().Set("Location", value)
	w.WriteHeader(302)
}

type ShortenRequestBody struct {
	URL string `json:"url"`
}
