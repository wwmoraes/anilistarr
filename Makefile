-include .env
-include .env.local
DOCKER ?= docker
GO ?= go
GORELEASER ?= goreleaser
export

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
coverage: coverage.out
	@${GO} tool cover -func=$<

%.out: $(SOURCE_FILES)
	@${GO} test -race -cover -coverprofile=$@ -v ./...

IMAGE ?= wwmoraes/anilistarr
# needs go install github.com/Khan/genqlient@latest
anilist:
	@cd internal/drivers/anilist && genqlient

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

run-coverage: GOCOVERDIR=coverage
run-coverage: PKGS="${PKG}/internal/usecases,${PKG}/internal/adapters"
run-coverage:
	@mkdir -p ${GOCOVERDIR}
	@echo "================================================================================"
	@echo "running application"
	@echo "================================================================================"
	${GO} run -cover -race -mod=readonly ./cmd/internal/integration/...
	@echo "================================================================================"
	@echo "raw coverage"
	@echo "================================================================================"
	${GO} tool covdata percent -i="${GOCOVERDIR}" | column -t
	@echo "================================================================================"
	@echo "coverage report"
	@echo "================================================================================"
	${GO} tool covdata textfmt -i="${GOCOVERDIR}" -o="${GOCOVERDIR}/gcov" -pkg="${PKGS}"
	${GO} tool cover -func="${GOCOVERDIR}/gcov"

run-integration: GITHUB_OUTPUT=tmp/github_output
run-integration: GOLANG_RUN_FLAGS=-race -mod=readonly
run-integration: GOLANG_PACKAGES=./cmd/internal/integration/...
run-integration: GOLANG_COVERAGE_PACKAGES=${PKG}/internal/usecases,${PKG}/internal/adapters
run-integration:
	@../actions/golang/integration/action.bash
