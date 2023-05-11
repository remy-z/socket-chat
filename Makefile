build:
	@go build -o bin/socket-chat

run: build 
	@./bin/socket-chat

test: 
	@go test -v ./...