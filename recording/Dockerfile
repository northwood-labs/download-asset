FROM alpine:latest@sha256:4bcff63911fcb4448bd4fdacec207030997caf25e9bea4045fa6c8c44de311d1 as fontcollector
RUN apk add --no-cache \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    font-meslo-nerd

FROM golang:alpine@sha256:daae04ebad0c21149979cd8e9db38f565ecefd8547cf4a591240dc1972cf1399 AS golang
COPY . /workspace
WORKDIR /workspace
RUN go install ./...

FROM ghcr.io/charmbracelet/vhs:latest@sha256:e88ed3faa06183a197fd44ded83e706098d9e4038b72da94bcdb9cb9b67e3527 AS debian
SHELL ["/bin/bash", "-euo", "pipefail", "-c"]
RUN apt-get -y update && apt-get install -y --no-install-recommends gpg wget
RUN mkdir -p /etc/apt/keyrings
RUN wget -qO- https://raw.githubusercontent.com/eza-community/eza/main/deb.asc | gpg --dearmor -o /etc/apt/keyrings/gierens.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/gierens.gpg] http://deb.gierens.de stable main" | tee /etc/apt/sources.list.d/gierens.list
RUN chmod 0644 /etc/apt/keyrings/gierens.gpg /etc/apt/sources.list.d/gierens.list
RUN apt-get -y update && apt-get install -y --no-install-recommends eza

FROM ghcr.io/charmbracelet/vhs:latest@sha256:e88ed3faa06183a197fd44ded83e706098d9e4038b72da94bcdb9cb9b67e3527 AS vhs
RUN apt-get -y update && apt-get install -y --no-install-recommends less && apt-get clean && rm -rf /var/lib/apt/lists/*
COPY --from=golang /go/bin/download-asset /usr/local/bin/download-asset
COPY --from=debian /usr/bin/eza /usr/local/bin/eza
COPY --from=fontcollector /usr/share/fonts/ /usr/share/fonts

ENTRYPOINT ["vhs"]
