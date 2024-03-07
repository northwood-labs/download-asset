.PHONY: build-setup
## build-setup: [build] Sets up the multi-arch build configuration.
build-setup:
	docker buildx use multiarch || docker buildx create --name multiarch --use 2>/dev/null

.PHONY: test
## test: [testing] Run tests.
test: build-setup
	DOCKER_BUILDKIT=1 docker buildx build \
		--output=type=docker \
		--load \
		--tag download-asset-test:latest \
		--secret id=GITHUB_TOKEN,env=GITHUB_TOKEN \
		--compress \
		--force-rm \
		--no-cache \
		--file bats/Dockerfile \
		.
