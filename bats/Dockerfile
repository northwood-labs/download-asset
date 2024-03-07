# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM golang:1.22-alpine AS go-installer

COPY download-asset.toml /etc/download-asset/download-asset.toml

COPY ./*.go /workspace/
COPY ./go.* /workspace/
COPY ./cmd/ /workspace/cmd/
COPY ./github/ /workspace/github/

WORKDIR /workspace

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo aquasecurity/trivy --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo direnv/direnv --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo golangci/golangci-lint --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo gruntwork-io/terragrunt --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo hadolint/hadolint --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo infracost/infracost --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo koalaman/shellcheck --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo iann0036/iamlive --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo minamijoyo/hcledit --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo minamijoyo/tfschema --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo opentofu/opentofu --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo terraform-docs/terraform-docs --verbose

RUN --mount=type=secret,id=GITHUB_TOKEN GITHUB_TOKEN="$(cat /run/secrets/GITHUB_TOKEN)" \
    go run main.go get --owner-repo terraform-linters/tflint --verbose

# ------------------------------------------------------------------------------

# syntax=docker/dockerfile:1
FROM --platform=$TARGETPLATFORM ubuntu:22.04 AS final
RUN apt-get -y update && apt-get -y upgrade && apt-get -y install file --no-install-recommends && rm -rf /var/lib/apt/lists/*

COPY --from=go-installer /usr/local/bin/trivy /usr/local/bin/trivy
COPY --from=go-installer /usr/local/bin/direnv /usr/local/bin/direnv
COPY --from=go-installer /usr/local/bin/golangci-lint /usr/local/bin/golangci-lint
COPY --from=go-installer /usr/local/bin/terragrunt /usr/local/bin/terragrunt
COPY --from=go-installer /usr/local/bin/hadolint /usr/local/bin/hadolint
COPY --from=go-installer /usr/local/bin/infracost /usr/local/bin/infracost
COPY --from=go-installer /usr/local/bin/shellcheck /usr/local/bin/shellcheck
COPY --from=go-installer /usr/local/bin/iamlive /usr/local/bin/iamlive
COPY --from=go-installer /usr/local/bin/hcledit /usr/local/bin/hcledit
COPY --from=go-installer /usr/local/bin/tfschema /usr/local/bin/tfschema
COPY --from=go-installer /usr/local/bin/tofu /usr/local/bin/tofu
COPY --from=go-installer /usr/local/bin/terraform-docs /usr/local/bin/terraform-docs
COPY --from=go-installer /usr/local/bin/tflint /usr/local/bin/tflint
