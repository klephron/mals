package uri_test

import (
	"mals-engine/pkg/uri"
	"testing"
)

func TestUriFile(t *testing.T) {
	source := "file:///home/nikit/tmp/main.py"
	expected := "/home/nikit/tmp/main.py"

	if actual, err := uri.UriToPath(source); err == nil {
		if actual != expected {
			t.Fatalf("expected %s, actual %s", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUriDirectory(t *testing.T) {
	source := "file:///home/nikit/tmp"
	expected := "/home/nikit/tmp"

	if actual, err := uri.UriToPath(source); err == nil {
		if actual != expected {
			t.Fatalf("expected %s, actual %s", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}

func TestUriDirectorySlash(t *testing.T) {
	source := "file:///home/nikit/tmp/"
	expected := "/home/nikit/tmp"

	if actual, err := uri.UriToPath(source); err == nil {
		if actual != expected {
			t.Fatalf("expected %s, actual %s", expected, actual)
		}
	} else {
		t.Fatal(err)
	}
}
