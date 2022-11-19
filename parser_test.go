package jreader

import (
	"net/http"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	DefaultFetchFunc = func(url string) (resp *http.Response, err error) {
		// Create mock http.Response
		resp = new(http.Response)
		resp.Body, err = os.Open("testdata/" + url)

		return resp, err
	}

	urls := []string{
		"rss_2.0.xml",
		"rss_2.0_content_encoded.xml",
		"rss_2.0_enclosure.xml",
		"rss_2.0-1.xml",
		"rss_2.0-1_enclosure.xml",
		"rss_1.0.xml",
		"atom_1.0.xml",
		"atom_1.0_enclosure.xml",
		"atom_1.0-1.xml",
		"atom_1.0_html.xml",
	}

	result, err := Parse(urls)
	if err != nil {
		t.Fatalf("Failed fetching testdata 'rss_2.0.xml': %v", err)
	}

	want := 92
	if len(result) != want {
		t.Errorf("got %d, want %d", len(result), want)
	}
}
