
define DOCKER_FLAGS
--platform=linux/arm64,linux/amd64
--label=org.opencontainers.image.created=${DATE}
--label=org.opencontainers.image.documentation=https://github.com/wwmoraes/anilistarr/blob/${BRANCH}/README.md
--label=org.opencontainers.image.revision=${REVISION}
--label=org.opencontainers.image.source=https://github.com/wwmoraes/anilistarr
--label=org.opencontainers.image.url=https://hub.docker.com/r/wwmoraes/anilistarr
--label=org.opencontainers.image.version=${VERSION}
--tag=wwmoraes/anilistarr:latest
--tag=wwmoraes/anilistarr:$(subst /,-,$(if ${BRANCH},${BRANCH},${REVISION}))
--tag=wwmoraes/anilistarr:$(subst +,-,${VERSION})
--tag=wwmoraes/anilistarr:$(shell echo ${VERSION} | cut -d- -f1)
--tag=wwmoraes/anilistarr:$(shell echo ${VERSION} | cut -d. -f1-2)
--tag=wwmoraes/anilistarr:$(shell echo ${VERSION} | cut -d. -f1)
endef

.PHONY: image
image: Dockerfile.tar

.PHONY: image-push
image-push: Dockerfile.tar
	## TODO use skopeo/regctl/oras to copy tarball instead
	docker buildx build --push $(filter-out --platform=%,$(strip ${DOCKER_FLAGS})) --file Dockerfile .

Dockerfile.tar: DOCKER_BUILDKIT=1
## avoids mixing application debugging with buildkit
Dockerfile.tar: GRPC_GO_LOG_VERBOSITY_LEVEL =
Dockerfile.tar: GRPC_GO_LOG_SEVERITY_LEVEL =
## https://github.com/moby/moby/issues/46129
Dockerfile.tar: OTEL_EXPORTER_OTLP_ENDPOINT =
Dockerfile.tar: Dockerfile dist
	docker buildx build $(strip ${DOCKER_FLAGS}) --output type=oci,dest=$@ --file $< .
	## do it again to load a single-architecture image into docker
	docker buildx build --load $(filter-out --platform=%,$(strip ${DOCKER_FLAGS})) --file $< .
	container-structure-test test -c container-structure-test.yaml -i wwmoraes/anilistarr:latest

Dockerfile.sarif: Dockerfile.hadolint.sarif Dockerfile.grype.sarif
	jq --slurp 'def deepmerge(a;b): reduce b[] as $$item (a; reduce ($$item | keys_unsorted[]) as $$key (.; $$item[$$key] as $$val | ($$val | type) as $$type | .[$$key] = if ($$type == "object") then deepmerge({}; [if .[$$key] == null then {} else .[$$key] end, $$val]) elif ($$type == "array") then (.[$$key] + $$val | unique) else $$val end)); deepmerge({}; .)' $^ > $@

Dockerfile.hadolint.sarif: Dockerfile
	-hadolint -f json $< | hadolint-sarif | tee $@ | sarif-fmt
	@jq -e '[.runs[].results[] | select(.level == "error")] | length | . == 0' $@ > /dev/null

Dockerfile.grype.sarif:
	grype db update || grype db delete && grype db update
	-grype -o sarif --fail-on critical wwmoraes/anilistarr:latest > $@
