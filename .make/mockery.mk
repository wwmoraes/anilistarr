define MOCKERY_SOURCES
$(strip $(shell git ls-files ':(attr:generates mockery)'))
endef

define MOCKERY_TARGETS
$(strip $(shell git ls-files ':(attr:generated mockery)'))
endef

${MOCKERY_TARGETS} &: ${MOCKERY_SOURCES}
	$(info updating mocks...)
	@go tool mockery
