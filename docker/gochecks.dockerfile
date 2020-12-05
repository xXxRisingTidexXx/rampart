FROM golangci/golangci-lint:v1.32.2-alpine

WORKDIR /go/src/app

ENV GOOS "linux"
ENV GOARCH "amd64"

COPY go.mod .
COPY go.sum .
COPY .golangci.yml .
COPY cmd cmd
COPY internal internal
COPY config config

RUN apk add build-base && go get ./... && go get github.com/kyoh86/richgo
