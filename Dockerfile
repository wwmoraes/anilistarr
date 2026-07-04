# syntax = docker/dockerfile:1

ARG ALPINE_VERSION=3

FROM alpine:${ALPINE_VERSION} AS certificates

# we always want the latest CA bundle
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates


FROM alpine:${ALPINE_VERSION} AS tmp

RUN mkdir -m 0750 /var/handler && chown 20000:20000 /var/handler


FROM alpine:${ALPINE_VERSION} AS machine-id

RUN sysctl -n kernel.random.uuid > /etc/machine-id


FROM scratch

WORKDIR /

COPY --from=machine-id /etc/machine-id /etc/machine-id

COPY --from=tmp --chown=20000:20000 --chmod=750 /var/handler /var/handler

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ARG TARGETOS
ARG TARGETARCH
COPY dist/handler_${TARGETOS}_${TARGETARCH}*/bin/handler /usr/bin/handler

USER 20000:20000

CMD ["/usr/bin/handler"]

ENV DATA_PATH=/var/handler
EXPOSE 8080
VOLUME /var/handler

LABEL org.opencontainers.image.authors="William Artero <docker@artero.dev>"
LABEL org.opencontainers.image.description="anilist custom list provider for sonarr/radarr"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/wwmoraes/anilistarr"
LABEL org.opencontainers.image.title="Anilistarr"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/wwmoraes/anilistarr"
LABEL org.opencontainers.image.vendor="William Artero <docker@artero.dev>"
