# mals

## Components

- [`mals`](https://github.com/klephron/mals.git) - main server component
- [`mals-adapter`](https://github.com/klephron/mals-adapter.git) - LSP stdio to TCP adapter
- [`mals-ctl`](https://github.com/klephron/mals-ctl.git) - utility command line tool
- [`mals-vscode`](https://github.com/klephron/mals-vscode.git) - VSCode Extension LSP client wrapper
- [`mals-test`](https://github.com/klephron/mals-test.git) - testing application and other testing utilities

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

Generate openapi schema (v3.1, v3.0):

```sh
curl $LISTENER_API_HOST:$LISTENER_API_PORT/openapi.{json,yml}
curl $LISTENER_API_HOST:$LISTENER_API_PORT/openapi-3.0.{json,yml}
```

## TODO

- Optimization of requests and caching:
  - Drop autocompletion requests that are old (may be checked by document version id) to skip answering to them if model can't pace
  - Skip if user triggers autocompletion too often and drop older requests
