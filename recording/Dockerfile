FROM alpine:latest@sha256:c5b1261d6d3e43071626931fc004f70149baeba2c8ec672bd4f27761f8e1ad6b as fontcollector
RUN apk add --no-cache \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    font-meslo-nerd

FROM golang:alpine@sha256:cdc86d9f363e8786845bea2040312b4efa321b828acdeb26f393faa864d887b0 AS golang
COPY . /workspace
WORKDIR /workspace
RUN go install ./...

FROM ghcr.io/charmbracelet/vhs:latest@sha256:87e8afac2837a79791af6dafdacb3047a49bb852acda371afc17ef3528d2a9b3 AS debian
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]
RUN apt-get -y update && apt-get install -y --no-install-recommends gpg wget
RUN mkdir -p /etc/apt/keyrings
RUN wget -qO- https://raw.githubusercontent.com/eza-community/eza/main/deb.asc | gpg --dearmor -o /etc/apt/keyrings/gierens.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/gierens.gpg] http://deb.gierens.de stable main" | tee /etc/apt/sources.list.d/gierens.list
RUN chmod 0644 /etc/apt/keyrings/gierens.gpg /etc/apt/sources.list.d/gierens.list
RUN apt-get -y update && apt-get install -y --no-install-recommends eza

FROM ghcr.io/charmbracelet/vhs:latest@sha256:87e8afac2837a79791af6dafdacb3047a49bb852acda371afc17ef3528d2a9b3 AS vhs
RUN apt-get -y update && apt-get install -y --no-install-recommends less && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY --from=golang /go/bin/download-asset /usr/local/bin/download-asset
COPY --from=debian /usr/bin/eza /usr/local/bin/eza
COPY --from=fontcollector /usr/share/fonts/ /usr/share/fonts

ENTRYPOINT ["vhs"]
