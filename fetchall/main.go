//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Stat struct {
	dur  time.Duration
	size int
	url  string
}

func fetch(url string, ch chan<- Stat) {

	now := time.Now()
	response, err := http.Get(url)
	if err != nil {
		ch <- (Stat{dur: time.Since(now), size: 0, url: url})
		return
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	ch <- (Stat{dur: time.Since(now), size: len(body), url: url})

}

func main() {
	st := make(chan Stat)

	for _, url := range os.Args[1:] {
		go fetch(url, st)
	}
	for range os.Args[1:] {
		s := <-st

		fmt.Printf("%v\t%v\t%v\n", s.dur, s.size, s.url)
	}
}
