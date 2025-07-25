# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.24-alpine@sha256:daae04ebad0c21149979cd8e9db38f565ecefd8547cf4a591240dc1972cf1399 AS base

COPY download-asset.toml /etc/download-asset/download-asset.toml

COPY ./*.go /workspace/
COPY ./go.* /workspace/
COPY ./cmd/ /workspace/cmd/
COPY ./github/ /workspace/github/

WORKDIR /workspace

# Compile the Go code.
RUN go run main.go --help

##
# And now… PARALLELISM!
##

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS trivy

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo aquasecurity/trivy

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS direnv

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo direnv/direnv

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS golangci-lint

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo golangci/golangci-lint

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS terragrunt

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo gruntwork-io/terragrunt

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS hadolint

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo hadolint/hadolint

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS infracost

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo infracost/infracost

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS shellcheck

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo koalaman/shellcheck

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS iamlive

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo iann0036/iamlive

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS hcledit

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo minamijoyo/hcledit

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS tfschema

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo minamijoyo/tfschema

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS opentofu

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo opentofu/opentofu

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS terraform-docs

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo terraform-docs/terraform-docs

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM base AS tflint

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo terraform-linters/tflint

#-------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM ubuntu:24.04@sha256:99c35190e22d294cdace2783ac55effc69d32896daaa265f0bbedbcde4fbe3e5 AS final
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install file --no-install-recommends && rm -rf /var/lib/apt/lists/*

COPY --from=trivy /usr/local/bin/trivy /usr/local/bin/trivy
COPY --from=direnv /usr/local/bin/direnv /usr/local/bin/direnv
COPY --from=golangci-lint /usr/local/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=terragrunt /usr/local/bin/terragrunt /usr/local/bin/terragrunt
COPY --from=hadolint /usr/local/bin/hadolint /usr/local/bin/hadolint
COPY --from=infracost /usr/local/bin/infracost /usr/local/bin/infracost
COPY --from=shellcheck /usr/local/bin/shellcheck /usr/local/bin/shellcheck
COPY --from=iamlive /usr/local/bin/iamlive /usr/local/bin/iamlive
COPY --from=hcledit /usr/local/bin/hcledit /usr/local/bin/hcledit
COPY --from=tfschema /usr/local/bin/tfschema /usr/local/bin/tfschema
COPY --from=opentofu /usr/local/bin/tofu /usr/local/bin/tofu
COPY --from=terraform-docs /usr/local/bin/terraform-docs /usr/local/bin/terraform-docs
COPY --from=tflint /usr/local/bin/tflint /usr/local/bin/tflint
