FROM scratch

WORKDIR /go/src/app
COPY . .

ENTRYPOINT ["/dist/datree_linux_amd64/datree"]
