FROM golang:1.15-alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -v ./...
RUN go install -v ./...

ENTRYPOINT ["/go/bin/datree"]