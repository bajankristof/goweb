FROM golang:1.26 AS toolchain

RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


FROM golang:1.26 AS base

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .


FROM base AS development

COPY --from=toolchain /go/bin/air /usr/local/bin/air
COPY --from=toolchain /go/bin/goose /usr/local/bin/goose
COPY --from=toolchain /go/bin/sqlc /usr/local/bin/sqlc

CMD ["air"]


FROM base AS build

RUN go build -o bin/goweb cmd/goweb


FROM debian:stable-slim

WORKDIR /app
COPY --from=build /app/bin/goweb /usr/local/bin/goweb

CMD ["goweb"]
