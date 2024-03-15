.PHONY: build-setup
## build-setup: [build] Sets up the multi-arch build configuration.
build-setup:
	docker buildx use multiarch || docker buildx create --name multiarch --use 2>/dev/null

.PHONY: test
## test: [testing] Run tests.
test: build-setup
	docker buildx build \
		--output=type=docker \
		--load \
		--tag download-asset-test:latest \
		--secret id=GITHUB_TOKEN,env=GITHUB_TOKEN \
		--compress \
		--force-rm \
		--no-cache \
		--file bats/Dockerfile \
		.

.PHONY: vhs-build
## vhs-build: [demo] Build the custom vhs image.
vhs-build:
	docker buildx build \
		--output=type=docker \
		--load \
		--tag vhs:latest \
		--no-cache \
		--compress \
		--force-rm \
		--file recording/Dockerfile \
		.

.PHONY: vhs
## vhs: [demo] Generate a demo video.
vhs:
	docker run --rm -e GITHUB_TOKEN -v "$$PWD:/vhs" vhs:latest recording/os-arch.tape
	docker run --rm -e GITHUB_TOKEN -v "$$PWD:/vhs" vhs:latest recording/get.tape
