# mals

## Generate LSP protocol structures

Build utility script to fetch latest changes:

```sh
go build -o build/lsp-gen cmd/lsp-gen/*.go
```

Generate lsp protocol structures and helpers:

```sh
./build/lsp-gen
```

## Build

```sh
go build -o build/mals cmd/mals/*.go
```

## Usage

```sh
./build/mals -h
```

## TODO

- Optimization of requests and caching:
  - Drop autocompletion requests that are old (may be checked by document version id) to skip answering to them if model can't pace
  - Skip if user triggers autocompletion too often and drop older requests
