-include .env
-include .env.local

GOCOVERDIR=coverage/integration
GOFLAGS += -cover -covermode=atomic -race -short -shuffle=on -mod=readonly -trimpath

export

.DEFAULT_GOAL := all
BRANCH != git branch --show-current
DATE != date --utc +"%Y-%m-%dT%H:%M:%SZ"
REVISION != git rev-parse --verify HEAD
VERSION != git describe --tags --always | cut -dv -f2

-include .make/*.mk

.PHONY: all
all: bin/handler gomod2nix.toml

.PHONY: configure
configure:
	cog install-hook --all --overwrite

.PHONY: check
check:
	@./.git/hooks/pre-commit
	@./.git/hooks/pre-push

.PHONY: clean
clean:
	-rm -rf bin build coverage dist
	-rm -f golangci-lint-report.xml *.sarif Dockerfile*.tar

.PHONY: dist
dist: dist/

.PHONY: test
test: coverage/all.txt

.PHONY: coverage
coverage: coverage/all.txt
	go tool cover -func=$< | sed 's|${GOMODULE}/||g'

.PHONY: release
release:
	goreleaser release --clean

.PHONY: generate
generate: $(strip ${GO_GENERATE_TARGETS})

.PHONY: coverage-upload
coverage-upload: coverage-unit-upload coverage-integration-upload

.PHONY: coverage-unit-upload
coverage-unit-upload: coverage/unit.part.txt coverage/unit.part.junit.xml
	codecov do-upload --disable-search --name unit --flag unit --file $(word 1,$^)
	codecov do-upload --disable-search --name unit --flag unit --file $(word 2,$^) --report-type test_results

.PHONY: coverage-integration-upload
coverage-integration-upload: coverage/integration.part.txt coverage/integration.part.junit.xml
	codecov do-upload --disable-search --name integration --flag integration --file $(word 1,$^)
	codecov do-upload --disable-search --name integration --flag integration --file $(word 2,$^) --report-type test_results

.PHONY: diagrams
diagrams: $(wildcard docs/structurizr-*.png) docs/components.png

.PHONY: fix
fix: fix-markdown
fix: fix-golang
fix: ;

.PHONY: fix-markdown
fix-markdown: $(shell git ls-files '*.md')
	$(info fixing in-place markdown files...)
	@markdownlint --fix $^

.PHONY: fix-golang
fix-golang:
	$(info fixing in-place golang issues...)
	@env -u GOFLAGS golangci-lint run --fix

.PHONY: sast
sast: sast.sarif

.PHONY: invoke-get-user
invoke-get-user:
	curl -v "http://${HOST}:${PORT}/user/${USERNAME}/id"

.PHONY: invoke-get-media
invoke-get-media:
	curl -v "http://${HOST}:${PORT}/user/${USERNAME}/media"

sast.sarif: Dockerfile.sarif semgrep.sarif
	jq --slurp 'def deepmerge(a;b): reduce b[] as $$item (a; reduce ($$item | keys_unsorted[]) as $$key (.; $$item[$$key] as $$val | ($$val | type) as $$type | .[$$key] = if ($$type == "object") then deepmerge({}; [if .[$$key] == null then {} else .[$$key] end, $$val]) elif ($$type == "array") then (.[$$key] + $$val | unique) else $$val end)); deepmerge({}; .)' $^ > $@

$(wildcard docs/structurizr-*.puml) &: docs/workspace.dsl
	structurizr-cli export -f plantuml/c4plantuml -w docs/workspace.dsl -o docs
	@touch docs/structurizr-*.puml

$(wildcard docs/structurizr-*.png) &: $(wildcard docs/structurizr-*.puml)
	plantuml $^
	@touch docs/structurizr-*.png

dist/: GOFLAGS=
dist/: ${GO_SOURCES} go.sum .goreleaser.yml
	goreleaser release --clean --snapshot --skip before

semgrep.sarif: ${GO_SOURCES} Dockerfile
	-semgrep scan --sarif-output=$@

docs/components.png: docs/components.puml
	plantuml $<
	@touch $@

%/:
	@test -d $@ || mkdir $@
