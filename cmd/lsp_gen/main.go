package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

const (
	source = "https://raw.githubusercontent.com/golang/tools/refs/heads/master"
)

func main() {
	files := []string{
		"gopls/internal/protocol/tsprotocol.go",
		"gopls/internal/protocol/tsdocument_changes.go",
	}

	params := argParse()

	os.MkdirAll(params.Target, os.ModePerm)

	for _, file := range files {
		sourcePath, err := url.JoinPath(source, file)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := http.Get(sourcePath)
		if err != nil {
			log.Fatal(err)
		}

		if resp.StatusCode != 200 {
			log.Fatalf("status code: %d %s", resp.StatusCode, file)
		}

		bytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("failed reading body: %s", err)
		}

		re := regexp.MustCompile(`(?m)^package\s+\w+`)
		modified := re.ReplaceAll(bytes, []byte("package "+params.Package))

		targetName := params.Prefix + filepath.Base(file)
		targetPath := path.Join(path.Join(params.Target, targetName))

		file, err := os.Create(targetPath)
		if err != nil {
			log.Fatalf("failed to create file: %s", targetPath)
		}

		_, err = file.Write(modified)
		if err != nil {
			log.Fatalf("%s", err.Error())
		}
	}
}
