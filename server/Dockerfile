FROM golang:latest AS base-build

WORKDIR /build
COPY go.mod .

RUN GOARCH=amd64 GOOS=linux go mod download

FROM base-build AS build

WORKDIR /build
COPY . .

RUN GOARCH=amd64 GOOS=linux go build -o filedrop cmd/file-drop/main.go

FROM scratch
WORKDIR /app
COPY --from=/build/filedrop /build/filedrop .
COPY ./server/template/dbSettings.json ./server/template/dbSettings.json

ENTRYPOINT ["./filedrop"]
