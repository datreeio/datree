FROM alpine:3.14 AS builder

RUN apk add --no-cache curl bash unzip openssl git

RUN curl https://get.datree.io | /bin/bash
RUN curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | /bin/bash
RUN helm plugin install https://github.com/datreeio/helm-datree
RUN mkdir /bin/plugintemp && cp -r $HOME/.local/share/helm/plugins/ /bin/plugintemp

RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl

FROM alpine:3.14
RUN apk add --no-cache bash

COPY --from=builder /usr/local/bin/datree /usr/local/bin/datree
COPY --from=builder /usr/local/bin/helm /usr/local/bin/helm
RUN mkdir /bin/plugintemp/
COPY --from=builder /bin/plugintemp/ /bin/plugintemp/
COPY --from=builder /usr/local/bin/kubectl /usr/local/bin/kubectl

RUN mkdir -p /root/.local/share/helm/ && cp -r /bin/plugintemp/plugins/ /root/.local/share/helm/

ENV HELM_PLUGINS="/root/.local/share/helm/plugins/"

ENTRYPOINT ["datree"]