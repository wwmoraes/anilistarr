# syntax = docker/dockerfile:1

ARG GOLANG_VERSION=1.22
ARG ALPINE_VERSION=3.20

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum tools.go ./
RUN go mod download -x

COPY *.go .
COPY cmd cmd
COPY pkg pkg
COPY internal internal
ARG VERSION
RUN --mount=type=cache,target=/root/.cache/go-build \
  go generate ./... && go build -o ./bin/handler ./cmd/handler/...


FROM alpine:${ALPINE_VERSION} AS certificates

# we always want the latest CA bundle
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates


FROM alpine:${ALPINE_VERSION} AS tmp
FROM scratch

LABEL org.opencontainers.image.authors="William Artero <docker@artero.dev>"
LABEL org.opencontainers.image.base.name="scratch"
LABEL org.opencontainers.image.description="anilist custom list provider for sonarr/radarr"
LABEL org.opencontainers.image.documentation="https://github.com/wwmoraes/anilistarr/blob/master/README.md"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/wwmoraes/anilistarr"
LABEL org.opencontainers.image.title="Anilistarr"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/wwmoraes/anilistarr"
LABEL org.opencontainers.image.vendor="William Artero <docker@artero.dev>"

CMD ["/usr/local/bin/handler"]
EXPOSE 8080

COPY --from=tmp --chown=20000:20000 /tmp /var/handler
COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /src/bin/handler /usr/local/bin/handler

WORKDIR /
USER 20000:20000
