# syntax=docker/dockerfile:1

# Build the application from source
FROM golang:1.22 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN cd ./cmd && go build -buildvcs=false -o /bot

# Deploy the application binary into a lean image
FROM debian:trixie-slim AS build-release-stage

WORKDIR /

# RUN apt-get -y update; apt-get -y install curl

COPY --from=build-stage /bot /bot

CMD ["/bot"]