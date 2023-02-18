FROM golang:1.17 AS base-build

WORKDIR /build
COPY go.mod .

RUN CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=linux go mod download

FROM base-build AS build

WORKDIR /build
COPY . .

RUN CGO_ENABLED=0 GOARCH=$TARGETARCH GOOS=linux go build -o cleaner cmd/cleaner/main.go

FROM scratch

WORKDIR /app
COPY --from=build /build/cleaner ./cleaner

ENTRYPOINT ["./cleaner"]
