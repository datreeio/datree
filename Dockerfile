FROM golang:1.15-alpine

WORKDIR /go/src/app
COPY . .

RUN go get -d -v ./...
RUN go build -tags main -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=0.9.0" -v
RUN go install -v ./...

ENTRYPOINT ["./datree"]
