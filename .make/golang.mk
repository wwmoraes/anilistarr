GO ?= go
COMMA := ,
SPACE := $(shell echo " ")

GOLANG_PACKAGE = $(shell ${GO} list)

GOCOVERDIR ?= coverage/integration

GOLANG_FLAGS ?= -race -mod=readonly
GOLANG_COVERAGE_PATH ?= coverage
GOLANG_OUTPUT_BIN_PATH ?= bin
GOLANG_RUN_REPORT_FILE ?= run-report.json

GOLANG_INTEGRATION_ENABLED ?=
GOLANG_INTEGRATION_SRC_PATH ?=
GOLANG_INTEGRATION_PACKAGES ?=

CMD_SOURCE_FILES := $(shell ${GO} list -f '{{ range .GoFiles }}{{ printf "%s/%s\n" $$.Dir . }}{{ end }}' ./cmd/...)
INTERNAL_SOURCE_FILES := $(shell ${GO} list -f '{{ range .GoFiles }}{{ printf "%s/%s\n" $$.Dir . }}{{ end }}' ./internal/...)
SOURCE_FILES := ${CMD_SOURCE_FILES} ${INTERNAL_SOURCE_FILES}

ifneq (${GOLANG_INTEGRATION_ENABLED},)
GOLANG_REPORT_SOURCE := ${GOLANG_COVERAGE_PATH}/merged.txt
else
GOLANG_REPORT_SOURCE := ${GOLANG_COVERAGE_PATH}/test.txt
endif

GOLANG_PKGS = $(subst ${SPACE},${COMMA},$(addprefix ${GOLANG_PACKAGE}/,$(subst ${COMMA},${SPACE},${GOLANG_INTEGRATION_PACKAGES})))

.PHONY: golang-build
golang-build:
	$(info generating files)
	@${GO} generate ./...
	$(info building binaries)
	@${GO} build ${GOLANG_FLAGS} -o ./${GOLANG_OUTPUT_BIN_PATH} ./...

.PHONY: golang-clean
golang-clean:
	-@${RM} -r ${GOLANG_COVERAGE_PATH}
	-@${RM} -r ${GOLANG_OUTPUT_BIN_PATH}

.PHONY: golang-test
golang-test:
	@${GO} test ${GOLANG_FLAGS} -v ./...

.PHONY: golang-report
golang-report: ${GOLANG_COVERAGE_PATH}/coverage.html

.PHONY: golang-coverage
golang-coverage: ${GOLANG_REPORT_SOURCE}
golang-coverage:
	$(info generating coverage report from $<)
	@${GO} tool cover -func="$<"

${GOLANG_COVERAGE_PATH}/coverage.html: ${GOLANG_REPORT_SOURCE}
	$(info generating html report from ${GOLANG_REPORT_SOURCE})
	@${GO} tool cover -html=$< -o $@

${GOLANG_COVERAGE_PATH}/test.txt: ${SOURCE_FILES}
${GOLANG_COVERAGE_PATH}/test.txt:
	$(info running unit tests)
	@${GO} test -v ${GOLANG_FLAGS} -coverprofile=$@ ./...

### only run and merge test results if we have a configured integration binary
ifneq (${GOLANG_INTEGRATION_ENABLED},)
${GOLANG_COVERAGE_PATH}/merged.txt: ${GOLANG_COVERAGE_PATH}/integration.txt ${GOLANG_COVERAGE_PATH}/test.txt
	$(info merging test results)
	@${GO} run github.com/wadey/gocovmerge@latest $^ > $@

${GOLANG_COVERAGE_PATH}/integration.txt: ${GOCOVERDIR}
	$(info converting integration results to gocov)
	@${GO} tool covdata textfmt -i="$<" -o="$@" -pkg="${GOLANG_PKGS}"

${GOCOVERDIR}: ${SOURCE_FILES}
${GOCOVERDIR}:
	$(info running integration test)
	@mkdir -p "$@"
	@${GO} tool test2json ${GO} run -cover ${GOLANG_FLAGS} ./${GOLANG_INTEGRATION_SRC_PATH}/... | tee ${GOLANG_RUN_REPORT_FILE}
	@${GO} tool covdata percent -i="$@" | column -t
endif
