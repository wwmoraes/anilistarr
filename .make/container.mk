define CONTAINER_NAMETAGS
$(strip
wwmoraes/anilistarr:latest
wwmoraes/anilistarr:$(subst /,-,$(if ${BRANCH},${BRANCH},${REVISION}))
wwmoraes/anilistarr:$(subst +,-,${VERSION})
wwmoraes/anilistarr:$(shell echo ${VERSION} | cut -d- -f1)
wwmoraes/anilistarr:$(shell echo ${VERSION} | cut -d. -f1-2)
wwmoraes/anilistarr:$(shell echo ${VERSION} | cut -d. -f1)
)
endef

define DOCKER_FLAGS
$(strip
--label=org.opencontainers.image.created=${DATE}
--label=org.opencontainers.image.documentation=https://github.com/wwmoraes/anilistarr/blob/${BRANCH}/README.md
--label=org.opencontainers.image.revision=${REVISION}
--label=org.opencontainers.image.version=${VERSION}
$(addprefix --tag=,${CONTAINER_NAMETAGS})
)
endef

ifneq ($(shell uname -s),Darwin)
DOCKER = docker buildx
DOCKER_FLAGS += --load
else
DOCKER = container
DOCKER_FLAGS += --platform=linux/arm64
endif

.PHONY: image
## avoids mixing application debugging with buildkit
image: GRPC_GO_LOG_VERBOSITY_LEVEL =
image: GRPC_GO_LOG_SEVERITY_LEVEL =
## https://github.com/moby/moby/issues/46129
image: OTEL_EXPORTER_OTLP_ENDPOINT =
image: Dockerfile dist/
	${DOCKER} build ${DOCKER_FLAGS} --file $< .
	container-structure-test test -c container-structure-test.yaml -i wwmoraes/anilistarr:latest

Dockerfile.sarif: Dockerfile.hadolint.sarif Dockerfile.grype.sarif
	$(info merging container SAST reports into $@...)
	@jq --slurp 'def deepmerge(a;b): reduce b[] as $$item (a; reduce ($$item | keys_unsorted[]) as $$key (.; $$item[$$key] as $$val | ($$val | type) as $$type | .[$$key] = if ($$type == "object") then deepmerge({}; [if .[$$key] == null then {} else .[$$key] end, $$val]) elif ($$type == "array") then (.[$$key] + $$val | unique) else $$val end)); deepmerge({}; .)' $^ > $@

Dockerfile.hadolint.sarif: Dockerfile
	$(info running SAST analysis (hadolint)...)
	-@hadolint -f json $< | hadolint-sarif | tee $@ | sarif-fmt
	@jq -e '[.runs[].results[] | select(.level == "error")] | length | . == 0' $@ > /dev/null

Dockerfile.grype.sarif: image
	$(info running SAST analysis (grype)...)
	-@grype db update --quiet || grype db delete --quiet && grype db update --quiet
	-@grype -o sarif --fail-on critical wwmoraes/anilistarr:latest | tee $@ | sarif-fmt
	@jq -e '[.runs[].results[] | select(.level == "error")] | length | . == 0' $@ > /dev/null
