CODECOV ?= codecov
CODECOV_FLAGS ?=
CODECOV_TOKEN ?=

.PHONY: codecov-report
codecov-report:
	@${CODECOV} create-report -t ${CODECOV_TOKEN} ${CODECOV_FLAGS}
