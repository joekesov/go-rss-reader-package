package jreader

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
	DateValid   bool
	Description string
}

func Parse(urls []string) ([]*RssItem, error) {
	rssItems := make([]*RssItem, 0)

	for _, stringUrl := range urls {
		feed, err := fetch(stringUrl)
		if err != nil {
			// handle error.
			return nil, err
		}

		fmt.Println(feed)
	}

	return rssItems, nil
}

// A FetchFunc is a function that fetches a feed for given URL.
type fetchFunc func(url string) (resp *http.Response, err error)

// DefaultFetchFunc uses http.DefaultClient to fetch a feed.
var DefaultFetchFunc = func(url string) (resp *http.Response, err error) {
	client := http.DefaultClient
	return client.Get(url)
}

// Parse RSS or Atom data.
func parse(data []byte) ([]*RssItem, error) {
	if strings.Contains(string(data), "<rss") {
		return parseRSS2(data)
	} else if strings.Contains(string(data), "xmlns=\"http://purl.org/rss/1.0/\"") {
		return parseRSS1(data)
	} else {
		return parseAtom(data)
	}
}

// Fetch downloads and parses the RSS feed at the given URL
func fetch(url string) ([]*RssItem, error) {
	return fetchByFunc(DefaultFetchFunc, url)
}

// FetchByFunc uses a func to fetch a URL.
func fetchByFunc(fetchFunc fetchFunc, url string) ([]*RssItem, error) {
	resp, err := fetchFunc(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	out, err := parse(body)
	if err != nil {
		return nil, err
	}

	return out, nil
}
