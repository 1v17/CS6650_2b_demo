# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.24-alpine  AS build
RUN apk add --no-cache git

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# disable cgo, target linux, static link
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o server .

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /app
COPY --from=build /src/server .

EXPOSE 8080
ENTRYPOINT ["./server"]