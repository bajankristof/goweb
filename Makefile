build: sqlc
	go build -o ./bin/goweb ./cmd/goweb

fmt:
	go fmt ./...

vet:
	go mod tidy
	go mod verify
	go vet ./...

lint:
	golangci-lint run

dev:
	air

sqlc:
	sqlc generate

migration:
	(cd sqlstore/postgresql/schema && goose create new sql)
	(cd sqlstore/sqlite/schema && goose create new sql)

# install: build
# 	@sudo ln -s $(shell pwd)/bin/goweb /usr/local/bin/goweb
