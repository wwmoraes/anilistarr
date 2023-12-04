GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOLANGCI_LINT_REPORT_FILE ?= golangci-lint-report.xml

GOLANGCI_LINT_SOURCE_FILES ?= $(shell ${GO} list -f '{{ range .GoFiles }}{{ printf "%s/%s\n" $$.Dir . }}{{ end }}' ./...)

.PHONY: golang-lint
golang-lint: ${GOLANGCI_LINT_REPORT_FILE}
	@${GOLANGCI_LINT} run

${GOLANGCI_LINT_REPORT_FILE}: ${GOLANGCI_LINT_SOURCE_FILES}
	@${GOLANGCI_LINT} run --out-format checkstyle > $@
