build:
	@swag init -g ./cmd/main.go -o ./docs --parseInternal true
	@go build -o bin/main cmd/main.go

run:build
	@bin/main