FROM ubuntu:18.04

WORKDIR /go/src/app
COPY . .

ENTRYPOINT ["./datree"]
