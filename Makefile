GIT_HASH ?= git-$(shell git rev-parse --short=12 HEAD)

IMAGE_NAME ?= quay.io/theauthgear/authgear-sms-gateway-tencent

.PHONY: build-image
build-image:
	# Add --pull so that we are using the latest base image.
	docker buildx build --platform linux/amd64 --pull --file ./Dockerfile --tag $(IMAGE_NAME):$(GIT_HASH) .

.PHONY: push-image
push-image: build-image
	docker push $(IMAGE_NAME):$(GIT_HASH)
