package jreader

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	urls := []string{
		"https://feeds.fireside.fm/bibleinayear/rss",
		"https://blog.centos.org/feed",
		"https://www.theregister.com/software/devops/headlines.atom",
	}

	got := Parse(urls)
	fmt.Println(got)
}

// TODO:
func TestAdd(t *testing.T) {
	a := 1
	b := 2
	expected := a + b

	if got := Add(a, b); got != expected {
		t.Errorf("Add(%d, %d) = %d, didn't return %d", a, b, got, expected)
	}
}

// TODO:
func TestSubtract(t *testing.T) {
	a := 1
	b := 2
	expected := a - b

	if got := Subtract(a, b); got != expected {
		t.Errorf("Subtract(%d, %d) = %d, didn't return %d", a, b, got, expected)
	}
}
