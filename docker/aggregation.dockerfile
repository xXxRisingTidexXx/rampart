# TODO: add distroless image (https://github.com/GoogleContainerTools/distroless).
# TODO: add non-root user.
FROM golang:1.15.3-alpine

WORKDIR /go/src/app

ARG RAMPART_UID
ARG RAMPART_GID

ENV GOOS "linux"
ENV GOARCH "amd64"

COPY go.mod .
COPY go.sum .
COPY cmd cmd
COPY internal internal
COPY config config

RUN apk add build-base && \
    go get ./... && \
    addgroup --gid $RAMPART_GID rampart && \
    adduser --disabled-password --gecos '' --uid $RAMPART_UID --ingroup rampart rampart
