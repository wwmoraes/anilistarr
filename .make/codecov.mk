CODECOV ?= codecov
CODECOV_FLAGS ?= -X fixes
CODECOV_TOKEN ?=

.PHONY: codecov-report
codecov-report:
	$(if $<,,$(error target codecov-report must have a source file as dependency))
	$(info uploading Codecov report)
	@${CODECOV} create-report -c -t ${CODECOV_TOKEN} ${CODECOV_FLAGS} -f "$<"
