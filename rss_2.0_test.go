package jreader

import (
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
