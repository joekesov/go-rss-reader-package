package jreader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
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
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()

	var errcList []<-chan error

	urlsc, errc, err := prepareUrls(ctx, urls...)
	if err != nil {
		return nil, err
	}
	errcList = append(errcList, errc)

	rssitemsc, errc, err := fetchItems(ctx, urlsc)
	if err != nil {
		return nil, err
	}
	errcList = append(errcList, errc)

	rssItems, errc, err := getItemsFromChan(ctx, rssitemsc)
	if err != nil {
		return nil, err
	}
	errcList = append(errcList, errc)

	err = waitForPipeline(errcList...)
	if err != nil {
		return nil, err
	}

	return rssItems, nil
}

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
func mergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	// We must ensure that the output channel has the capacity to hold as many errors
	// as there are error channels. This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))

	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	// Start a goroutine to close out once all the output goroutines are
	// done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

// WaitForPipeline waits for results from all error channels.
// It returns early on the first error.
func waitForPipeline(errs ...<-chan error) error {
	errc := mergeErrors(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

func prepareUrls(ctx context.Context, urls ...string) (<-chan string, <-chan error, error) {
	if len(urls) == 0 {
		// Handle an error that occurs before the goroutine begins.
		return nil, nil, fmt.Errorf("no lines provided")
	}
	out := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)

		for index, stringUrl := range urls {
			if stringUrl == "" {
				// Handle an error that occurs during the goroutine.
				errc <- fmt.Errorf("url %v is empty", index+1)
				return
			}

			select {
			case out <- stringUrl:
			case <-ctx.Done():
				return
			}
		}
	}()

	return out, errc, nil
}

func fetchItems(ctx context.Context, in <-chan string) (<-chan RssItem, <-chan error, error) {
	addedStream := make(chan RssItem)
	errc := make(chan error, 1)

	go func() {
		defer close(addedStream)
		defer close(errc)

		for stringUrl := range in {
			feed, err := fetch(stringUrl)
			if err != nil {
				errc <- fmt.Errorf("can't fetch for url:'%s': %v", stringUrl, err)
				return
			}

			for _, item := range feed {
				select {
				case addedStream <- *item:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return addedStream, errc, nil
}

func getItemsFromChan(ctx context.Context, in <-chan RssItem) ([]*RssItem, <-chan error, error) {
	rssItems := make([]*RssItem, 0)
	errc := make(chan error, 1)
	defer close(errc)

	for item := range in {
		rssItems = append(rssItems, &item)
	}

	return rssItems, errc, nil
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
