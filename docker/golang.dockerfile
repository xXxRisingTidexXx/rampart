FROM golang:1.15.3-alpine

WORKDIR /go/src/app

ENV GOOS "linux"
ENV GOARCH "amd64"

COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY internal internal
COPY config config
COPY templates templates

RUN apk add build-base && go get ./...