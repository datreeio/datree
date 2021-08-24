FROM ubuntu:18.04

WORKDIR /go/src/app
COPY . .
RUN apt-get update && apt-get install -y curl grep unzip bash
RUN curl https://get.datree.io | /bin/bash

RUN go get -d -v ./...
RUN go build -tags main -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.7.0-staging" -v
RUN go install -v ./...

ENTRYPOINT ["/go/bin/datree"]
