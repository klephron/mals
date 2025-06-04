package uri_test

import (
	"mals-engine/pkg/url"
	"testing"
)

func TestUrlFile(t *testing.T) {
	source := "file:///home/nikit/tmp/main.py"
	expected := "/home/nikit/tmp/main.py"

	if actual, err := url.UriToPath(source); err == nil {
		if actual != expected {
			t.Fatalf("expected %s, actual %s", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUrlDirectory(t *testing.T) {
	source := "file:///home/nikit/tmp"
	expected := "/home/nikit/tmp"

	if actual, err := url.UriToPath(source); err == nil {
		if actual != expected {
			t.Fatalf("expected %s, actual %s", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUrlDirectorySlash(t *testing.T) {
	source := "file:///home/nikit/tmp/"
	expected := "/home/nikit/tmp"

	if actual, err := url.UriToPath(source); err == nil {
		if actual != expected {
			t.Fatalf("expected %s, actual %s", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}
