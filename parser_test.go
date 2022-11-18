package jreader

import (
	"fmt"
	"net/http"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	fetch1 := makeTestdataFetchFunc("rss_2.0.xml")
	got, err := fetchByFunc(fetch1, "http://localhost/dummyrss")
	if err != nil {
		t.Fatalf("Failed fetching testdata 'rss_2.0.xml': %v", err)
	}

	for _, i := range got {
		fmt.Println(i)
	}
}

func makeTestdataFetchFunc(file string) fetchFunc {
	return func(url string) (resp *http.Response, err error) {
		// Create mock http.Response
		resp = new(http.Response)
		resp.Body, err = os.Open("testdata/" + file)

		return resp, err
	}
}
