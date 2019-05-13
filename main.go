package main

import (
	"encoding/xml"
	"log"
	"io/ioutil"
	"net/http"
	"time"
	"html/template"
	"path/filepath"
	"strings"
)

var feedURLs = []string{
	"https://blog.golang.org/feed.atom",
	"http://feeds.kottke.org/main",
}

func main() {
	feeds := fetchFeeds(feedURLs)
	feeds = filterFeeds(feeds)
	if len(feeds) == 0 {
		return
	}
	t := loadTemplate()
	writer := &strings.Builder{}
	err := t.Execute(writer, feeds)
	if err != nil {
		log.Fatalln(err)
	}
	err = mail(writer.String())
	if err != nil {
		log.Fatalln(err)
	}
}

func loadTemplate() *template.Template {
	templates, err := filepath.Glob("templates/*")
	if err != nil {
		log.Println(err)
	}
	t := template.Must(template.New("layout.go.html").ParseFiles(templates...))
	return t
}

func filterFeeds(feeds []*Feed) ([]*Feed) {
	cutoff := time.Now().Add(-24*time.Hour)
	filtered := make([]*Feed, 0, len(feeds))
	for _, f := range feeds {
		if time.Time(f.Updated).After(cutoff) {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

type feedResult struct {
	feed *Feed
	err error
}

func fetchFeeds(urls []string) ([]*Feed) {
	// Create a channel to process the feeds
	feedc := make(chan feedResult, len(urls))

	// Start a goroutine for each feed url
	for _, u := range urls {
		go fetchFeed(u, feedc)
	}

	// Wait for the goroutines to write their results to the channel
	feeds := []*Feed{}
	for i := 0; i < len(urls); i++ {
		res := <-feedc
		// If the goroutine errors out, we'll just wait for others
		if res.err != nil {
			continue
		}
		feeds = append(feeds, res.feed)
	}

	return feeds
}

func fetchFeed(url string, feedc chan feedResult) {
	// Create a client with a default timeout
	net := &http.Client{
		Timeout: time.Second * 10,
	}
	// Issue a GET request for the feed
	res, err := net.Get(url)
	// If there was an error write that to the channel and return immediately
	if err != nil {
		feedc <- feedResult{nil, err}
		return
	}
	defer res.Body.Close()
	// Read the body of the request and parse the feed
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		feedc <- feedResult{nil, err}
		return
	}
	feed, err := parseFeed(body)
	if err != nil {
		feedc <- feedResult{nil, err}
		return
	}
	feedc <- feedResult{feed, nil}
}

func parseFeed(body []byte) (*Feed, error) {
	feed := Feed{}
	err := xml.Unmarshal(body, &feed)
	if err != nil {
		return nil, err
	}
	return &feed, nil
}
