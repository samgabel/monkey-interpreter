.SILENT:

run:
	go run .

build:
	go build .

test:
	go test ./...

test-ast:
	go test -v ./ast

test-lexer:
	go test -v ./lexer

test-parser:
	go test -v ./parser

format:
	gofmt -w -d .
