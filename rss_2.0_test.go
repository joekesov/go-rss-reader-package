package jreader

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestParseItemLen(t *testing.T) {
	tests := map[string]int{
		"rss_2.0.xml":                 2,
		"rss_2.0_content_encoded.xml": 1,
		"rss_2.0_enclosure.xml":       1,
		"rss_2.0-1.xml":               4,
		"rss_2.0-1_enclosure.xml":     1,
	}

	for test, want := range tests {
		name := filepath.Join("testdata", test)
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("Reading %s: %v", name, err)
		}

		feed, err := parse(data)
		if err != nil {
			t.Fatalf("Parsing %s: %v", name, err)
		}

		if len(feed) != want {
			t.Errorf("%s: got %d, want %d", name, len(feed), want)
		}
	}
}

func TestParseItemDateOK(t *testing.T) {
	tests := map[string]string{
		"rss_2.0.xml":                 "2009-09-06 16:45:00 +0000 UTC",
		"rss_2.0_content_encoded.xml": "2009-09-06 16:45:00 +0000 UTC",
		"rss_2.0_enclosure.xml":       "2009-09-06 16:45:00 +0000 UTC",
		"rss_2.0-1.xml":               "2003-06-03 09:39:21 +0000 UTC",
		"rss_2.0-1_enclosure.xml":     "2016-05-14 15:39:34 +0000 UTC",
	}

	for test, want := range tests {
		name := filepath.Join("testdata", test)
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("Reading %s: %v", name, err)
		}

		feed, err := parse(data)
		if err != nil {
			t.Fatalf("Parsing %s: %v", name, err)
		}

		if !feed[0].DateValid {
			t.Errorf("%s: date %q invalid!", name, feed[0].PublishDate)
		} else if got := feed[0].PublishDate.UTC().String(); got != want {
			t.Errorf("%s: got %q, want %q", name, got, want)
		}
	}
}

func TestParseItemDateFailure(t *testing.T) {
	tests := map[string]string{
		"rss_2.0.xml": "0001-01-01 00:00:00 +0000 UTC",
	}

	for test, want := range tests {
		name := filepath.Join("testdata", test)
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("Reading %s: %v", name, err)
		}

		feed, err := parse(data)
		if err != nil {
			t.Fatalf("Parsing %s: %v", name, err)
		}

		if fmt.Sprintf("%s", feed[1].PublishDate) != want {
			t.Errorf("%s: got %q, want %q", name, feed[1].PublishDate, want)
		}

		if feed[1].DateValid {
			t.Errorf("%s: got unexpected valid date", name)
		}
	}
}
