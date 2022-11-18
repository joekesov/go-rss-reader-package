package jreader

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

func parseRSS1(data []byte) ([]*RssItem, error) {
	feed := rss1_0Feed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}
	if feed.Channel == nil {
		return nil, fmt.Errorf("no channel found in %q", string(data))
	}

	channel := feed.Channel

	out := make([]*RssItem, 0, len(feed.Items))
	// Process items.
	for _, item := range feed.Items {
		next := new(RssItem)
		next.Title = item.Title
		next.Source = channel.Title
		next.SourceURL = channel.Link
		next.Description = item.Description
		next.Link = item.Link
		if item.Date != "" {
			next.PublishDate, err = parseTime(item.Date)
			if err == nil {
				next.DateValid = true
			}
		} else if item.PubDate != "" {
			next.PublishDate, err = parseTime(item.PubDate)
			if err == nil {
				next.DateValid = true
			}
		}

		out = append(out, next)
	}

	return out, nil
}

type rss1_0Feed struct {
	XMLName xml.Name       `xml:"RDF"`
	Channel *rss1_0Channel `xml:"channel"`
	Items   []rss1_0Item   `xml:"item"`
}

type rss1_0Channel struct {
	XMLName xml.Name `xml:"channel"`
	Title   string   `xml:"title"`
	Link    string   `xml:"link"`
}

type rss1_0Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Date        string   `xml:"date"`
	DateValid   bool
}
