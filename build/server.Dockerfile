# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.21 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 cd ./cmd && go build -buildvcs=false -o /server

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM debian:trixie-slim AS build-release-stage

WORKDIR /

RUN apt-get -y update; apt-get -y install curl

COPY --from=build-stage /server /server

EXPOSE $HTTP_PORT

CMD ["/server"]