FROM ubuntu:18.04

WORKDIR /go/src/app
COPY . .

ENTRYPOINT ["bash", "/dist/datree_linux_amd64/datree"]
