package jreader

import (
	"time"
)

type RssItem struct {
	Title       string
	Source      string
	SourceURL   string
	Link        string
	PublishDate time.Time
}

func Parse(urls []string) []RssItem {
	var rssItems []RssItem

	return rssItems
}

// TODO:
// Returns the sum of two numbers
func Add(a int, b int) int {
	return a + b
}

// TODO:
// Returns the difference between two numbers
func Subtract(a int, b int) int {
	return a - b
}
