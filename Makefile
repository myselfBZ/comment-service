build:
	@go build -o bin/main ./cmd/api/

run: build 
	@./bin/main
	@go run ./cmd/consumer/main.go
