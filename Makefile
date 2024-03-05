# Disable built-in rules and variables and suffixes
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:

-include .env
-include .env.local
export

GOLANG_INTEGRATION_ENABLED = 1
CONTAINER_STRUCTURE_TEST_FILE = container-structure-test.yaml
CONTAINER_IMAGE = wwmoraes/anilistarr
GOLANG_INTEGRATION_SRC_PATH = cmd/internal/integration
GOLANG_INTEGRATION_PACKAGES = internal/usecases,internal/adapters

-include .make/*.mk

codecov-report: ${GOLANG_COVERAGE_REPORT_SOURCE}
${CC_REPORT_JSON_PATH}: ${GOLANG_COVERAGE_REPORT_SOURCE}

### local targets

.PHONY: build
build: golang-build

.PHONY: clean
clean: golang-clean

.PHONY: release
release: golang-release

.PHONY: test
test: golang-test

.PHONY: lint
lint: golang-lint container-lint

.PHONY: lint-report
lint-report: golang-lint-report container-lint-report

.PHONY: coverage
coverage: golang-coverage

.PHONY: coverage-html
coverage-html: golang-coverage-html

.PHONY: report
report: code-climate-report codecov-report

.PHONY: image
image: container-image

.PHONY: image-test
image-test: container-test

# needs go install github.com/Khan/genqlient@latest
anilist:
	@cd internal/drivers/anilist && genqlient

run:
	@${GO} run ./cmd/handler/...

redis-cli:
	@redis-cli -p 16379

redis-proxy:
	@flyctl redis proxy

get-user:
	@curl -v "http://${HOST}:${PORT}/user/wwmoraes/id"

get-media:
	@curl -v "http://${HOST}:${PORT}/user/wwmoraes/media"

docs:
	@open http://localhost:6060/pkg/github.com/${IMAGE}/
	@godoc -notes="BUG|TODO|FIX"

run-image: DATA_PATH=/var/handler
run-image:
	@${CONTAINER} run --rm \
	-e DATA_PATH \
	-it ${IMAGE}

diagrams:
	@structurizr-cli export -f plantuml/c4plantuml -w workspace.dsl -o docs
	@plantuml docs/*.puml

watch-diagrams:
	@fswatch -o --event Updated workspace.dsl | xargs -n 1 sh -c "clear; date; ${MAKE} diagrams"

# edit-diagrams:
# 	@ structurizr/lite

api:
	@mkdir -p internal/api
	@oapi-codegen \
		-generate types,chi-server,spec \
		-package api \
		-o internal/api/api.gen.go \
		swagger.yaml
	@sed -i '' '/var err error/d' internal/api/api.gen.go
