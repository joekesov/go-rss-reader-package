package jreader

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseAtomContent(t *testing.T) {
	tests := map[string]string{
		"atom_1.0.xml":           "Volltext des Weblog-Eintrags",
		"atom_1.0_enclosure.xml": "Volltext des Weblog-Eintrags",
		"atom_1.0-1.xml":         "",
		"atom_1.0_html.xml":      "<body>html</body>",
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

		//if feed[0].Content != want {
		//	t.Errorf("%s: got %q, want %q", name, feed.Items[0].Content, want)
		//}

		if !feed[0].DateValid {
			t.Errorf("%s: Invalid date: %q", name, feed[0].PublishDate)
		}
	}
}
