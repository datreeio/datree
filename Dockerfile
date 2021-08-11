FROM ubuntu:18.04

WORKDIR /go/src/app
COPY . .

# ENTRYPOINT ["./dist/datree_linux_amd64/datree"]
CMD ["ls", "-l"]
