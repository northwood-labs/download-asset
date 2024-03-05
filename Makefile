.PHONY: test
## test: [testing] Run tests.
test:
	DOCKER_BUILDKIT=1 docker buildx build \
		--output=type=docker \
		--load \
		--tag download-asset-test:latest \
		--secret id=GITHUB_TOKEN,env=GITHUB_TOKEN \
		--compress \
		--force-rm \
		--file bats/Dockerfile \
		.

# 		--no-cache \
