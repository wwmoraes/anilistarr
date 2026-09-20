GOCOVERDIR ?= coverage/integration

GOMODULE != go list -m

define GENERATE_TARGETS
$(strip $(shell git ls-files -- ':(attr:generated)'))
endef

define GO_TEST_IGNORE_PATTERNS
$(strip $(addsuffix :,$(shell git ls-files -- ':(attr:generated)')))
endef

define TEST_PACKAGES
$(strip $(shell git ls-files -- ':(attr:!skip-coverage)**.go' | xargs dirname | sort -u))
endef

define GO_SOURCES
$(strip $(shell git ls-files ':(attr:!generated !test)*.go'))
endef

IMPURE_OUTPUT = coverage/impure.unit.json
IMPURE_PROFILE = coverage/impure.unit.txt
PURE_OUTPUT = coverage/pure.unit.json
PURE_PROFILE = coverage/pure.unit.txt

check:: ${GENERATE_TARGETS}
	golangci-lint run

gomod2nix.toml: go.sum
	gomod2nix generate --with-deps
	gomod2nix import --with-deps
	@echo update default.nix with the new vendor hash: $(shell nix-hash --sri --type sha256 $(nix build --repair --no-link --print-out-paths .#handler.goModules))

go.sum: ${GO_SOURCES} go.mod
	@go mod tidy -v -x
	@touch $@

bin/handler: ${GO_SOURCES} go.sum
	go build -race -mod=readonly -trimpath -ldflags="-s -w -X 'main.version=${VERSION}-${REVISION}'" -o ./$@ ./cmd/handler/...

${PURE_PROFILE} ${PURE_OUTPUT} &: ${GO_SOURCES} ${MOCKERY_TARGETS} go.sum | coverage/
	go test -covermode=atomic -race -shuffle=on -mod=readonly -json -coverprofile=${PURE_PROFILE} -tags=pure $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES})) | tee ${PURE_OUTPUT} | gotestdox
	sed -i'' '\%$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))%d' $@

${IMPURE_PROFILE} ${IMPURE_OUTPUT} &: ${GO_SOURCES} ${MOCKERY_TARGETS} go.sum | coverage/
	go test -covermode=atomic -race -shuffle=on -mod=readonly -json -coverprofile=${IMPURE_PROFILE} $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES})) | tee ${IMPURE_OUTPUT} | gotestdox
	sed -i'' '\%$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))%d' $@

coverage/unit.part.txt: coverage/pure.unit.txt coverage/impure.unit.txt
	go tool gocovmerge $^ > $@

coverage/integration.part.txt: ${GO_SOURCES} go.sum | coverage/ ${GOCOVERDIR}/
	-@rm -rf "${GOCOVERDIR}/*" 2>/dev/null || true
	env GOCOVERDIR=${GOCOVERDIR} go run -cover -covermode=atomic -race -mod=readonly ./cmd/internal/integration/...
	go tool covdata textfmt -i=${GOCOVERDIR} -o=$@ -pkg="$(addprefix ${GOMODULE}/,${TEST_PACKAGES})"
	sed -i'' '\%$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))%d' $@

coverage/all.txt: coverage/unit.part.txt coverage/integration.part.txt
	go tool gocovmerge $^ \
	| grep $(if ${GO_TEST_IGNORE_PATTERNS},-v '$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS})',.) \
	> $@

coverage/%.junit.xml: coverage/%.txt
	go-junit-report -in $< -out $@

coverage/%.svg: coverage/%.txt
	go-cover-treemap -coverprofile $< > $@

coverage/%.html: coverage/%.txt
	go tool cover -html=$< -o $@
