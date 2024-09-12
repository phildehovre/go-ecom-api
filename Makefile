build:
	@go build -o bin/go-complete-api cmd/main.go

test:
	@go test -v ./...

run:
	@./bin/go-complete-api