FROM ubuntu:18.04

WORKDIR /app
RUN apt-get update && apt-get install -y curl grep unzip bash
RUN curl https://get.datree.io | /bin/bash

ENTRYPOINT ["datree"]
