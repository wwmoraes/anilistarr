-include .env
-include .env.local
export

CMD_SOURCE_FILES := $(shell find cmd -type f -name '*.go')
INTERNAL_SOURCE_FILES := $(shell find internal -type f -name '*.go')
SOURCE_FILES := $(CMD_SOURCE_FILES) $(INTERNAL_SOURCE_FILES)

.PHONY: build
build: generate
	go build -o ./bin/ ./...

.PHONY: release
release:
	goreleaser release --clean --skip-publish --skip-announce --snapshot

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test -race -v ./...

.PHONY: coverage
coverage: coverage.out
	@go tool cover -func=$<

%.out: $(SOURCE_FILES)
	@go test -race -cover -coverprofile=$@ -v ./...

IMAGE ?= wwmoraes/anilistarr
# needs go install github.com/Khan/genqlient@latest
anilist:
	@cd internal/drivers/anilist && genqlient

image: CREATED=$(shell date -u +"%Y-%m-%dT%TZ")
image: REVISION=$(shell git log -n 1 --format="%H")
image: VERSION=$(patsubst v%,%,$(shell git describe --tags 2> /dev/null || echo "0.1.0-rc.0"))
image:
	docker build --load $(if ${TARGET},--target ${TARGET}) \
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
	@go run ./cmd/handler/...

redis-cli:
	@redis-cli -p 16379

redis-proxy:
	@flyctl redis proxy

get-user:
	@curl -v "http://127.0.0.1:8080/user?name=wwmoraes"
