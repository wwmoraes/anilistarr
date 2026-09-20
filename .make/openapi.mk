GENERATE_TARGETS += internal/api/api.gen.go

internal/api/api.gen.go: swagger.yaml
	$(info updating OpenAPI bindings...)
	@go tool oapi-codegen -generate types,chi-server,spec -package api -o $@ $<
	# Thank you oapi-codegen...
	@sed --in-place= "/var err error/d" $@
	@sed --in-place= "/_ = err/d" $@
