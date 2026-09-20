define SQLC_SOURCES
$(strip $(shell git ls-files ':(attr:generates sqlc)'))
endef

define SQLC_TARGETS
$(strip $(shell git ls-files ':(attr:generated sqlc)'))
endef

${SQLC_TARGETS} &: ${SQLC_SOURCES}
	$(info updating SQL bindings...)
	@go tool sqlc generate
