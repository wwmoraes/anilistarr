GORELEASER ?= goreleaser

.PHONY: golang-release
golang-release:
	@${GORELEASER} release --clean --skip-publish --skip-announce --snapshot
