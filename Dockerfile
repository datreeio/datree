FROM scratch
ENTRYPOINT ["dist/datree-macos_darwin_amd64/datree"]
COPY datree /
