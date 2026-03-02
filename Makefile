build:
	@go build -o ./bin/watchbowl ./cmd/watchbowl

fmt:
	@go fmt ./...

sqlc:
	@sqlc generate

migration:
	@cd db/schema && goose create new sql
