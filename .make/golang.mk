GO ?= go
COMMA := ,
SPACE := $(shell echo " ")

GOLANG_PACKAGE = $(shell ${GO} list)

GOCOVERDIR ?= coverage/integration

GOLANG_FLAGS ?= -race -mod=readonly
GOLANG_OUTPUT_BIN_PATH ?= bin

GOLANG_COVERAGE_PATH ?= coverage
GOLANG_COVERAGE_REPORT_TARGET ?= ${GOLANG_COVERAGE_PATH}/coverage.html
GOLANG_INTEGRATION_TEST_GOCOV_FILE ?= ${GOLANG_COVERAGE_PATH}/integration.txt
GOLANG_MERGED_GOCOV_FILE ?= ${GOLANG_COVERAGE_PATH}/merged.txt
GOLANG_RUN_REPORT_FILE ?= ${GOLANG_COVERAGE_PATH}/run-report.json
GOLANG_UNIT_TEST_GOCOV_FILE ?= ${GOLANG_COVERAGE_PATH}/unit.txt

GOLANG_INTEGRATION_SRC_PATH ?=
GOLANG_INTEGRATION_PACKAGES ?=
GOLANG_INTEGRATION_ENABLED ?= $(if ${GOLANG_INTEGRATION_SRC_PATH},1)

CMD_SOURCE_FILES := $(shell ${GO} list -f '{{ range .GoFiles }}{{ printf "%s/%s\n" $$.Dir . }}{{ end }}' ./cmd/...)
INTERNAL_SOURCE_FILES := $(shell ${GO} list -f '{{ range .GoFiles }}{{ printf "%s/%s\n" $$.Dir . }}{{ end }}' ./internal/...)
SOURCE_FILES := ${CMD_SOURCE_FILES} ${INTERNAL_SOURCE_FILES}

### use either the unit test or the merged gocov file for reports
ifneq (${GOLANG_INTEGRATION_ENABLED},)
GOLANG_COVERAGE_REPORT_SOURCE := ${GOLANG_MERGED_GOCOV_FILE}
else
GOLANG_COVERAGE_REPORT_SOURCE := ${GOLANG_UNIT_TEST_GOCOV_FILE}
endif

### prefix the relative paths with the package name. Required by go tool covdata
GOLANG_PKGS = $(subst ${SPACE},${COMMA},$(addprefix ${GOLANG_PACKAGE}/,$(subst ${COMMA},${SPACE},${GOLANG_INTEGRATION_PACKAGES})))

.PHONY: golang-build
golang-build:
	$(info generating files)
	@${GO} generate ./...
	$(info building binaries)
	@${GO} build ${GOLANG_FLAGS} -o ./${GOLANG_OUTPUT_BIN_PATH}/ ./...

.PHONY: golang-clean
golang-clean:
	-@${RM} ${GOLANG_COVERAGE_REPORT_TARGET}
	-@${RM} ${GOLANG_INTEGRATION_TEST_GOCOV_FILE}
	-@${RM} ${GOLANG_MERGED_GOCOV_FILE}
	-@${RM} ${GOLANG_RUN_REPORT_FILE}
	-@${RM} ${GOLANG_UNIT_TEST_GOCOV_FILE}
	-@${RM} -r ${GOLANG_COVERAGE_PATH}
	-@${RM} -r ${GOLANG_OUTPUT_BIN_PATH}

.PHONY: golang-test
golang-test:
	@${GO} test ${GOLANG_FLAGS} -v ./...

.PHONY: golang-report
golang-report: ${GOLANG_COVERAGE_REPORT_TARGET}

.PHONY: golang-coverage
golang-coverage: ${GOLANG_COVERAGE_REPORT_SOURCE}
golang-coverage:
	$(info generating coverage report from $<)
	@${GO} tool cover -func="$<"

### make magic, not war :)

### gocov => HTML
${GOLANG_COVERAGE_REPORT_TARGET}: ${GOLANG_COVERAGE_REPORT_SOURCE}
	$(info generating html report from ${GOLANG_COVERAGE_REPORT_SOURCE})
	@${GO} tool cover -html=$< -o $@

### Go source => unit test gocov
${GOLANG_UNIT_TEST_GOCOV_FILE}: ${SOURCE_FILES}
${GOLANG_UNIT_TEST_GOCOV_FILE}:
	$(info running unit tests)
	@${GO} test -v ${GOLANG_FLAGS} -coverprofile=$@ ./...

### only run integration test and merge test results if enabled
ifneq (${GOLANG_INTEGRATION_ENABLED},)

### gocov files => merged gocov
${GOLANG_MERGED_GOCOV_FILE}: ${GOLANG_INTEGRATION_TEST_GOCOV_FILE}
${GOLANG_MERGED_GOCOV_FILE}: ${GOLANG_UNIT_TEST_GOCOV_FILE}
${GOLANG_MERGED_GOCOV_FILE}:
	$(info merging test results)
	@${GO} run github.com/wadey/gocovmerge@latest $^ > $@

### instrumentation data => gocov
${GOLANG_INTEGRATION_TEST_GOCOV_FILE}: ${GOCOVERDIR}
	$(info converting integration results to gocov)
	@${GO} tool covdata textfmt -i="$<" -o="$@" -pkg="${GOLANG_PKGS}"

### Go source => instrumentation data
${GOCOVERDIR}: ${SOURCE_FILES}
${GOCOVERDIR}:
	$(info running integration test)
	@mkdir -p "$@"
	@${GO} tool test2json ${GO} run -cover ${GOLANG_FLAGS} ./${GOLANG_INTEGRATION_SRC_PATH}/... | tee ${GOLANG_RUN_REPORT_FILE}
	@${GO} tool covdata percent -i="$@" | column -t
endif
