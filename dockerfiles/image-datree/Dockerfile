FROM golang:1.18-alpine AS builder
RUN apk --no-cache add curl

WORKDIR /go/src/app
COPY . .

RUN curl --silent "https://api.github.com/repos/datreeio/datree/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' > cli-version
RUN go get -d -v ./...
RUN go build -tags main -ldflags="-extldflags '-static' -X github.com/datreeio/datree/cmd.CliVersion=$(cat cli-version)" -v

FROM alpine:3.14
COPY --from=builder /go/src/app/datree /
ENTRYPOINT ["/datree"]
