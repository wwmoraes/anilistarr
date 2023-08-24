-include .env
-include .env.local
DOCKER ?= docker
GO ?= go
GORELEASER ?= goreleaser
GOCOVERDIR ?= coverage/integration
export

SEPARATOR = $(shell printf "%0.s=" {1..80})

PKG = github.com/wwmoraes/anilistarr
CMD_SOURCE_FILES := $(shell find cmd -type f -name '*.go')
INTERNAL_SOURCE_FILES := $(shell find internal -type f -name '*.go')
SOURCE_FILES := $(CMD_SOURCE_FILES) $(INTERNAL_SOURCE_FILES)

.PHONY: build
build: generate
	$(info building)
	@${GO} build -o ./bin/ ./...

.PHONY: release
release:
	@${GORELEASER} release --clean --skip-publish --skip-announce --snapshot

.PHONY: generate
generate:
	$(info generating files)
	@${GO} generate ./...

.PHONY: test
test:
	@${GO} test -race -v ./...

.PHONY: coverage
coverage: GOCOVERDIR=coverage/integration
coverage: PKGS="${PKG}/internal/usecases,${PKG}/internal/adapters"
coverage: coverage/coverage.txt
coverage:
	$(info ${SEPARATOR})
	$(info coverage report)
	$(info ${SEPARATOR})
	@${GO} tool cover -func="$<"

PHONY: coverage-html
coverage-html: coverage/coverage.html

${GOCOVERDIR}: ${SOURCE_FILES}
${GOCOVERDIR}:
	$(info ${SEPARATOR})
	$(info running integration test)
	$(info ${SEPARATOR})
	@mkdir -p "$@"
	@${GO} run -cover -race -mod=readonly ./cmd/internal/integration/...
	@echo "${SEPARATOR}"
	@echo "raw report"
	@echo "${SEPARATOR}"
	@${GO} tool covdata percent -i="$@" | column -t

%/coverage.txt: PKGS="${PKG}/internal/usecases,${PKG}/internal/adapters"
%/coverage.txt: ${GOCOVERDIR}
	$(info ${SEPARATOR})
	$(info generating gcov data)
	$(info ${SEPARATOR})
	@${GO} tool covdata textfmt -i="$<" -o="$@" -pkg="${PKGS}"

%/coverage.html: %/coverage.txt
	$(info ${SEPARATOR})
	$(info generating html report)
	$(info ${SEPARATOR})
	@${GO} tool cover -html=$< -o $@

IMAGE ?= wwmoraes/anilistarr
# needs go install github.com/Khan/genqlient@latest
anilist:
	@cd internal/drivers/anilist && genqlient

## https://github.com/moby/moby/issues/46129
image: OTEL_EXPORTER_OTLP_ENDPOINT=
image: CREATED=$(shell date -u +"%Y-%m-%dT%TZ")
image: REVISION=$(shell git log -n 1 --format="%H")
image: VERSION=$(patsubst v%,%,$(shell git describe --tags 2> /dev/null || echo "0.1.0-rc.0"))
image:
	$(info building image ${IMAGE})
	@${DOCKER} build --load $(if ${TARGET},--target ${TARGET}) \
		-t ${IMAGE} \
		--build-arg VERSION=${VERSION} \
		--label org.opencontainers.image.created=${CREATED} \
		--label org.opencontainers.image.revision=${REVISION} \
		--label org.opencontainers.image.documentation=https://github.com/${IMAGE}/blob/master/README.md \
		--label org.opencontainers.image.source=https://github.com/${IMAGE} \
		--label org.opencontainers.image.url=https://hub.docker.com/r/${IMAGE} \
		.
	@container-structure-test test -c container-structure-test.yaml -i wwmoraes/anilistarr

run:
	@${GO} run ./cmd/handler/...

redis-cli:
	@redis-cli -p 16379

redis-proxy:
	@flyctl redis proxy

get-user:
	@curl -v "http://127.0.0.1:8080/user?name=wwmoraes"
