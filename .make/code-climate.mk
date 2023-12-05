CC_TEST_REPORTER ?= cc-test-reporter

CC_REPORT_JSON_PATH ?= coverage/codeclimate.json

code-climate-report: ${CC_REPORT_JSON_PATH}
	$(info uploading report to Code Climate)
	@${CC_TEST_REPORTER} upload-coverage

${CC_REPORT_JSON_PATH}: ${GOLANG_COVERAGE_REPORT_SOURCE}
	$(info converting gocov to Code Climate format)
	@${CC_TEST_REPORTER} format-coverage -o $@ -t gocov -p ${GOLANG_PACKAGE} $<
