FROM golang:1.14.4-alpine

WORKDIR /go/src/app

ENV GOOS "linux"
ENV GOARCH "amd64"
ENV CGO_ENABLED 0

COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY internal internal
COPY config config

RUN go get ./...
