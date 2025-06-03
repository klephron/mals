package uri

import (
	"net/url"
	"path/filepath"
)

func UriToPath(rawURL string) (string, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", err
	}

	path := filepath.Clean(u.Path)

	return path, nil
}
