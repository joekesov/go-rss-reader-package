package jreader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseRSS(t *testing.T) {
	tests := map[string]string{
		"rss_1.0.xml": "",
	}

	for test, _ := range tests {
		name := filepath.Join("testdata", test)
		data, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("Reading %s: %v", name, err)
		}

		feed, err := parse(data)
		if err != nil {
			t.Fatalf("Parsing %s: %v", name, err)
		}

		if len(feed) != 40 {
			t.Errorf("%v: expected 40 items, got: %v", name, len(feed))
		} else {
			for i, item := range feed {
				if !item.DateValid {
					t.Errorf("%v Invalid date for item (#%v): %v", name, i, item.Title)
				}
			}
		}
	}
}
