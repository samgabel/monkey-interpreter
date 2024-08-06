.SILENT:

run:
	go run .

build:
	go build .

test:
	go test ./...

test-lexer:
	go test -v ./lexer

format:
	gofmt -w -d .
