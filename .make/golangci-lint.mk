GO ?= go
GOLANGCI_LINT ?= golangci-lint
GOLANGCI_LINT_REPORT_FILE ?= golangci-lint-report.json
GOLANGCI_LINT_FORMAT ?= json

GOLANGCI_LINT_SOURCE_FILES ?= $(shell ${GO} list -f '{{ range .GoFiles }}{{ printf "%s/%s\n" $$.Dir . }}{{ end }}' ./...)

.PHONY: golang-lint
golang-lint:
	$(info linting go source)
	@${GOLANGCI_LINT} run

golang-lint-report: ${GOLANGCI_LINT_REPORT_FILE}

${GOLANGCI_LINT_REPORT_FILE}: ${GOLANGCI_LINT_SOURCE_FILES}
	$(info generating lint report of go source)
	@${GOLANGCI_LINT} run --out-format json > $@
