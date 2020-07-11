FROM golangci/golangci-lint:v1.28.2-alpine

WORKDIR /go/src/app

COPY go.mod .
COPY go.sum .
COPY .golangci.yml .
COPY cmd cmd
COPY internal internal
COPY config config
COPY secrets secrets

RUN go get -t ./...
