VERSION := $(shell cat VERSION)
COMMIT := $(shell git rev-parse --short HEAD)
LONG_VERSION := $(VERSION)-$(COMMIT)

.PHONY: build
build:
	docker build -t dagan/rustler -t dagan/rustler:v$(VERSION) -t dagan/rustler:v$(LONG_VERSION) .

.PHONY: clean
clean: IMAGE_ID=$(shell docker image ls -q dagan/rustler:v$(LONG_VERSION))
clean: RAIDER_IMAGE_ID=$(shell docker image ls -q dagan/rustler:raider-v$(LONG_VERSION))
clean:
	[ -z "$(IMAGE_ID)" ] || docker image rm --force $(IMAGE_ID)
	[ -z "$(RAIDER_IMAGE_ID)" ] || docker image rm --force $(RAIDER_IMAGE_ID)

.PHONY: cobra
cobra:
	docker build -t dagan/rustler-cobra --target cobra-cli .

.PHONY: raider
raider: build
	docker build -f docker/raider.dockerfile -t dagan/rustler:raider -t dagan/rustler:raider-v$(VERSION) -t dagan/rustler:raider-v$(LONG_VERSION) .