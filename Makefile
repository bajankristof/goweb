build:
	@go build -o ./bin/goweb ./cmd/goweb

fmt:
	@go fmt ./...

sqlc:
	@sqlc generate

migration:
	@cd db/schema && goose create new sql
