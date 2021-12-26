FROM golang:latest AS base-build

WORKDIR /build
COPY go.mod .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go mod download

FROM base-build AS build

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o filedrop-cleaner cmd/cleaner/main.go

FROM scratch
WORKDIR /app
COPY ./db_settings/dbSettings-docker.json ./db_settings/dbSettings.json
COPY --from=build /build/filedrop-cleaner ./filedrop-cleaner

ENTRYPOINT ["./filedrop-cleaner"]
