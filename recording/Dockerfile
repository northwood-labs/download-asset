FROM alpine:latest@sha256:beefdbd8a1da6d2915566fde36db9db0b524eb737fc57cd1367effd16dc0d06d as fontcollector
RUN apk add --no-cache \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/main \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/community \
    --repository=http://dl-cdn.alpinelinux.org/alpine/edge/testing \
    font-meslo-nerd

FROM golang:alpine@sha256:2a882244fb51835ebbd8313bffee83775b0c076aaf56b497b43d8a4c72db65e1 AS golang
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
