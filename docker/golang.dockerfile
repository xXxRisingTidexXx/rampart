FROM golang:1.16.3-alpine3.13

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