package jreader

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

func parseRSS2(data []byte) ([]*RssItem, error) {
	f := rss2_0Feed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	err := p.Decode(&f)
	if err != nil {
		return nil, err
	}
	if f.Channel == nil {
		return nil, fmt.Errorf("no channel found in %q", string(data))
	}

	channel := f.Channel
	channelLink := ""
	for _, link := range channel.Link {
		if link.Rel == "" && link.Type == "" && link.Href == "" && link.Chardata != "" {
			channelLink = link.Chardata
			break
		}
	}

	out := make([]*RssItem, 0, len(channel.Items))
	// Process items.
	for _, i := range channel.Items {
		next := new(RssItem)
		next.Title = i.Title
		next.Source = channel.Title
		next.SourceURL = channelLink
		next.Link = i.Link
		if i.PubDate != "" {
			next.PublishDate, err = parseTime(i.PubDate)
			if err == nil {
				next.DateValid = true
			}
		} else if i.Date != "" {
			next.PublishDate, err = parseTime(i.Date)
			if err == nil {
				next.DateValid = true
			}
		}

		next.Description = i.Description

		out = append(out, next)
	}

	return out, nil
}

type rss2_0Feed struct {
	XMLName xml.Name       `xml:"rss"`
	Channel *rss2_0Channel `xml:"channel"`
}

type rss2_0Channel struct {
	XMLName     xml.Name     `xml:"channel"`
	Title       string       `xml:"title"`
	Description string       `xml:"description"`
	Link        []rss2_0Link `xml:"link"`
	Items       []rss2_0Item `xml:"item"`
}

type rss2_0Link struct {
	Rel      string `xml:"rel,attr"`
	Href     string `xml:"href,attr"`
	Type     string `xml:"type,attr"`
	Chardata string `xml:",chardata"`
}

type rss2_0Item struct {
	XMLName     xml.Name `xml:"item"`
	Title       string   `xml:"title"`
	Description string   `xml:"description"`
	Link        string   `xml:"link"`
	PubDate     string   `xml:"pubDate"`
	Date        string   `xml:"date"`
	DateValid   bool
	ID          string `xml:"guid"`
}
