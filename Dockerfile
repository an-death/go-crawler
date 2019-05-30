MAINTAINER 'as'

FROM golang:1.11

COPY . /crawler
WORKDIR /crawler

RUN go build -mod vendor -a  -ldflags \
    "-s -w \
    -X main.VERSION=0.0.1 -X main.BUILD=$(date -u +%Y-%m-%d/%H:%M:%S) \
    "

