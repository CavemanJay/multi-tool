FROM golang AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN ping 8.8.8.8
# RUN go get ./...
# ARG VERSION
# RUN go build -ldflags="-s -w -X main.version=${VERSION}" -o gogurt