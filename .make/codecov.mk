CODECOV ?= codecov
CODECOV_FLAGS ?= -X fixes -f ${GOLANG_COVERAGE_REPORT_SOURCE}
CODECOV_TOKEN ?=

.PHONY: codecov-report
codecov-report: ${GOLANG_COVERAGE_REPORT_SOURCE}
	$(info uploading Codecov report)
	@${CODECOV} create-report -t ${CODECOV_TOKEN} ${CODECOV_FLAGS}
