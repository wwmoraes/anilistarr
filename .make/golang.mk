GOCOVERDIR ?= coverage/integration

GOMODULE != go list -m

define GO_GENERATE_SOURCES
$(strip
db/queries.sql
internal/drivers/sqlite/schema.sql
sqlc.yaml
swagger.yaml
)
endef

define GO_GENERATE_TARGETS
$(strip
internal/api/api.gen.go
internal/drivers/sqlite/model/db.go
internal/drivers/sqlite/model/models.go
internal/drivers/sqlite/model/queries.sql.go
)
endef

define GO_TEST_IGNORE_PATTERNS
$(strip
.gen.go:
.pb.go:
.pb.gw.go:
.sql.go:
.xo.go:
)
endef

define TEST_PACKAGES
$(strip
internal/adapters
internal/api
internal/drivers
internal/entities
internal/usecases
pkg
)
endef

define MOCKERY_SOURCES
$(strip
$(shell git ls-files ':(exclude)**_test.go' 'internal/usecases/*.go')
)
endef

define MOCKERY_TARGETS
$(strip
$(wildcard internal/test/*_mocks.go)
)
endef

GO_SOURCES = $(filter-out ${MOCKERY_TARGETS},$(shell git ls-files '*.go') $(strip ${GO_GENERATE_TARGETS}))

IMPURE_OUTPUT = coverage/impure.unit.json
IMPURE_PROFILE = coverage/impure.unit.txt
PURE_OUTPUT = coverage/pure.unit.json
PURE_PROFILE = coverage/pure.unit.txt

check:: ${GO_GENERATE_TARGETS}
	golangci-lint run

gomod2nix.toml: go.sum
	gomod2nix generate

go.sum: ${GO_SOURCES} go.mod
	@go mod tidy -v -x
	@touch $@

bin/handler: ${GO_SOURCES} go.sum
	go build -race -mod=readonly -trimpath -ldflags="-s -w -X 'main.version=${VERSION}-${REVISION}'" -o ./$@ ./cmd/handler/...

${GO_GENERATE_TARGETS} &: ${GO_GENERATE_SOURCES}
	go generate ./...

${PURE_PROFILE} ${PURE_OUTPUT} &: ${GO_SOURCES} ${MOCKERY_TARGETS} go.sum | coverage/
	go test -covermode=atomic -race -shuffle=on -mod=readonly -json -coverprofile=${PURE_PROFILE} -tags=pure $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES})) | tee ${PURE_OUTPUT} | gotestdox
	sed -i'' '/$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))/d' $@

${IMPURE_PROFILE} ${IMPURE_OUTPUT} &: ${GO_SOURCES} ${MOCKERY_TARGETS} go.sum | coverage/
	go test -covermode=atomic -race -shuffle=on -mod=readonly -json -coverprofile=${IMPURE_PROFILE} $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES})) | tee ${IMPURE_OUTPUT} | gotestdox
	sed -i'' '/$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))/d' $@

coverage/unit.part.txt: coverage/pure.unit.txt coverage/impure.unit.txt
	go run github.com/wadey/gocovmerge $^ > $@

coverage/integration.part.txt: ${GO_SOURCES} go.sum | coverage/ ${GOCOVERDIR}/
	-@rm -rf "${GOCOVERDIR}/*" 2>/dev/null || true
	env GOCOVERDIR=${GOCOVERDIR} go run -cover -covermode=atomic -race -mod=readonly ./cmd/internal/integration/...
	go tool covdata textfmt -i=${GOCOVERDIR} -o=$@ -pkg="$(addprefix ${GOMODULE}/,${TEST_PACKAGES})"
	sed -i'' '/$(subst .,\.,$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS}))/d' $@

coverage/all.txt: coverage/unit.part.txt coverage/integration.part.txt
	go run github.com/wadey/gocovmerge $^ \
	| grep $(if ${GO_TEST_IGNORE_PATTERNS},-v '$(subst $(space),\|,${GO_TEST_IGNORE_PATTERNS})',.) \
	> $@

coverage/%.junit.xml: coverage/%.txt
	go-junit-report -in $< -out $@

coverage/%.svg: coverage/%.txt
	go-cover-treemap -coverprofile $< > $@

coverage/%.html: coverage/%.txt
	go tool cover -html=$< -o $@

${MOCKERY_TARGETS} &: .mockery.yml ${MOCKERY_SOURCES}
	$(info updating mocks...)
	@go tool mockery
