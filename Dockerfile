FROM ubuntu:18.04

WORKDIR /go/src/app
COPY . .

CMD ["sh", "-c", "tail -f /dev/null"]
