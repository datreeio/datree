FROM golang:1.15-alpine
RUN apk --no-cache add curl

WORKDIR /go/src/app
COPY . .

RUN curl --silent "https://api.github.com/repos/datreeio/datree/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/' > cli-version
RUN go get -d -v ./...
RUN go build -tags main -ldflags="-X github.com/datreeio/datree/cmd.CliVersion=$(cat cli-version)" -v
RUN go install -v ./...

ENTRYPOINT ["./datree"]
