GOCOVERDIR ?= coverage/integration
GOFLAGS += -cover -covermode=atomic -race -short -shuffle=on -mod=readonly -trimpath

export GOCOVERDIR GOFLAGS

GOMODULE != go list -m

define GO_GENERATE_SOURCES
db/queries.sql
internal/drivers/sqlite/schema.sql
sqlc.yaml
swagger.yaml
endef

define GO_GENERATE_TARGETS
internal/api/api.gen.go
internal/drivers/sqlite/model/db.go
internal/drivers/sqlite/model/models.go
internal/drivers/sqlite/model/queries.sql.go
endef

define TEST_PACKAGES
internal/adapters
internal/api
internal/drivers
internal/entities
internal/usecases
pkg
endef

GO_SOURCES = $(shell git ls-files '*.go') $(strip ${GO_GENERATE_TARGETS})

go.sum: GOFLAGS-=-mod-readonly
go.sum: ${GO_SOURCES} go.mod
	@go mod tidy -v -x
	@touch $@

bin/handler: ${GO_SOURCES} go.sum
	go build -ldflags="-s -w -X 'main.version=${VERSION}-${REVISION}'" -o ./$@ ./cmd/handler/...

$(strip ${GO_GENERATE_TARGETS}) &: $(strip ${GO_GENERATE_SOURCES})
	env -u GOCOVERDIR go generate ./...

coverage/pure.unit.txt: ${GO_SOURCES} go.sum | coverage/
	go test -coverprofile=$@ -tags=pure $(strip $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES})))
	sed -i'' '#\.gen\.go:|\.pb\.go:|\.pb\.gw\.go:|\.sql\.go:|\.xo\.go:#d' $@

coverage/impure.unit.txt: ${GO_SOURCES} go.sum | coverage/
	go test -coverprofile=$@ $(strip $(addprefix ./,$(addsuffix /...,${TEST_PACKAGES})))
	sed -i'' '#\.gen\.go:|\.pb\.go:|\.pb\.gw\.go:|\.sql\.go:|\.xo\.go:#d' $@

coverage/unit.part.txt: coverage/pure.unit.txt coverage/impure.unit.txt
	go run github.com/wadey/gocovmerge $^ > $@

coverage/integration.part.txt: ${GO_SOURCES} go.sum | coverage/ ${GOCOVERDIR}/
	-@rm -rf "${GOCOVERDIR}/*" 2>/dev/null || true
	go run ./cmd/internal/integration/...
	go tool covdata textfmt -i=${GOCOVERDIR} -o=$@ -pkg="$(strip $(addprefix ${GOMODULE}/,${TEST_PACKAGES}))"
	sed -i'' '#\.gen\.go:|\.pb\.go:|\.pb\.gw\.go:|\.sql\.go:|\.xo\.go:#d' $@

coverage/all.txt: coverage/unit.part.txt coverage/integration.part.txt
	go run github.com/wadey/gocovmerge $^ > $@

coverage/%.junit.xml: coverage/%.txt
	go-junit-report -in $< -out $@

coverage/%.svg: coverage/%.txt
	go-cover-treemap -coverprofile $< > $@

coverage/%.html: coverage/%.txt
	go tool cover -html=$< -o $@
