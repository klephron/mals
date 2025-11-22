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

func httpRead(source string, file string) []byte {
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

	return bytes
}

func write(target string, bytes []byte) {
	dir := path.Dir(target)
	os.MkdirAll(dir, os.ModePerm)
	file, err := os.Create(target)
	if err != nil {
		log.Fatalf("failed to create file: %s", target)
	}
	_, err = file.Write(bytes)
	if err != nil {
		log.Fatalf("%s", err.Error())
	}
}

func modify(params Params, bytes []byte) []byte {
	re := regexp.MustCompile(`(?m)^package\s+\w+`)
	modified := re.ReplaceAll(bytes, []byte("package "+params.Package))
	return modified
}

func getTargetPath(params Params, file string) string {
	targetName := params.Prefix + filepath.Base(file)
	targetPath := path.Join(path.Join(params.Target, targetName))
	return targetPath
}

func genGopls(params Params) {
	files := []string{
		"gopls/internal/protocol/tsprotocol.go",
		"gopls/internal/protocol/tsdocument_changes.go",
	}

	for _, file := range files {
		bytes := httpRead(source, file)
		modified := modify(params, bytes)
		targetPath := getTargetPath(params, file)
		write(targetPath, modified)
	}
}

func genUri(params Params) {
	bytes := []byte(`package protocol

type URI = string
type DocumentURI = string`)

	modified := modify(params, bytes)
	targetPath := getTargetPath(params, "uri.go")
	write(targetPath, modified)
}

func main() {
	params := argParse()

	genGopls(params)
	genUri(params)
}
