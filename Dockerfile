FROM ubuntu:18.04

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -tags main -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.8.0-staging" -v
RUN go install -v ./...

ENTRYPOINT ["/go/bin/datree"]
