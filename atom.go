package jreader

import (
	"bytes"
	"encoding/xml"
)

func parseAtom(data []byte) ([]*RssItem, error) {
	feed := atomFeed{}
	p := xml.NewDecoder(bytes.NewReader(data))
	err := p.Decode(&feed)
	if err != nil {
		return nil, err
	}

	sourceUrl := ""
	for _, link := range feed.Link {
		if link.Rel == "alternate" || link.Rel == "" {
			sourceUrl = link.Href
			break
		}
	}

	out := make([]*RssItem, 0, len(feed.Items))

	// Process items.
	for _, item := range feed.Items {
		next := new(RssItem)
		next.Title = item.Title
		next.Source = feed.Title
		next.SourceURL = sourceUrl
		next.Description = item.Summary
		if item.Date != "" {
			next.PublishDate, err = parseTime(item.Date)
			if err == nil {
				next.DateValid = true
			}
		}
		for _, link := range item.Links {
			if link.Rel == "alternate" || link.Rel == "" {
				next.Link = link.Href
				break
			}
		}

		out = append(out, next)
	}

	return out, nil
}

type RAWContent struct {
	RAWContent string `xml:",innerxml"`
}

type atomFeed struct {
	XMLName     xml.Name   `xml:"feed"`
	Title       string     `xml:"title"`
	Description string     `xml:"subtitle"`
	Link        []atomLink `xml:"link"`
	Items       []atomItem `xml:"entry"`
}

type atomItem struct {
	XMLName   xml.Name   `xml:"entry"`
	Title     string     `xml:"title"`
	Summary   string     `xml:"summary"`
	Content   RAWContent `xml:"content"`
	Links     []atomLink `xml:"link"`
	Date      string     `xml:"updated"`
	DateValid bool
	ID        string `xml:"id"`
}

type atomLink struct {
	Href   string `xml:"href,attr"`
	Rel    string `xml:"rel,attr"`
	Type   string `xml:"type,attr"`
	Length uint   `xml:"length,attr"`
}
