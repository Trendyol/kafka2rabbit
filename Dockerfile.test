FROM oms-registry.trendyol.com:5005/base/golang:1.13.1-alpine3.10 AS builder
ARG DOCKER_HOST
ARG DOCKER_DRIVER
ENV DOCKER_HOST=$DOCKER_HOST
ENV DOCKER_DRIVER=$DOCKER_DRIVER
#ENV TARGET_FOLDER=rabbitmq

RUN apk update && apk add git curl gcc musl-dev
RUN apk add --update --no-cache ca-certificates && rm -rf /var/cache/apk/*
RUN go get github.com/onsi/ginkgo/ginkgo
RUN go get github.com/onsi/gomega/...

WORKDIR /go/src/github.com/kafka2rabbit
ADD . .
RUN go mod download
RUN go build

ENTRYPOINT ginkgo -v -r services/event_executor
