-include .env
-include .env.local
export

.PHONY: build
build:
	go generate ./...
	go build -o ./bin/ ./...

IMAGE ?= wwmoraes/anilistarr
# needs go install github.com/Khan/genqlient@latest
anilist:
	@cd internal/anilist && genqlient

image:
	docker build --load $(if ${TARGET},--target ${TARGET}) \
		-t ${IMAGE} \
		.

run:
	@go run ./cmd/handler/...

redis-cli:
	@redis-cli -p 16379

redis-proxy:
	@flyctl redis proxy

get-user:
	@curl -v "http://127.0.0.1:8080/user?name=wwmoraes"
