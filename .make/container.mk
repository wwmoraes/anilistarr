CONTAINER ?= docker
CONTAINER_STRUCTURE_TEST ?= container-structure-test

CONTAINER_IMAGE ?=
CONTAINER_STRUCTURE_TEST_FILE ?=

## https://github.com/moby/moby/issues/46129
container-image: OTEL_EXPORTER_OTLP_ENDPOINT=
container-image: CREATED=$(shell date -u +"%Y-%m-%dT%TZ")
container-image: REVISION=$(shell git log -n 1 --format="%H")
container-image: VERSION=$(patsubst v%,%,$(shell git describe --tags 2> /dev/null || echo "0.1.0-rc.0"))
container-image:
	$(info building image ${CONTAINER_IMAGE})
	@${CONTAINER} build --load $(if ${TARGET},--target ${TARGET}) \
		-t ${CONTAINER_IMAGE} \
		--build-arg VERSION=${VERSION} \
		--label org.opencontainers.image.created=${CREATED} \
		--label org.opencontainers.image.revision=${REVISION} \
		--label org.opencontainers.image.version=${VERSION} \
		.

container-test: ${CONTAINER_STRUCTURE_TEST_FILE}
	@${CONTAINER_STRUCTURE_TEST} test -c "$<" -i "${CONTAINER_IMAGE}"
