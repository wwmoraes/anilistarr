# syntax = docker/dockerfile:1

ARG GOLANG_VERSION=1.23
ARG ALPINE_VERSION=3.21
ARG DATA_PATH=/var/handler

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS build

WORKDIR /src

COPY go.mod go.sum tools.go ./
RUN go mod download -x

COPY *.go .
COPY cmd cmd
COPY pkg pkg
COPY internal internal
COPY db db
COPY sqlc.yaml swagger.yaml ./
RUN go generate ./...

ARG VERSION
RUN --mount=type=cache,target=/root/.cache/go-build \
  go build -o ./bin/handler ./cmd/handler/...


FROM alpine:${ALPINE_VERSION} AS certificates

# we always want the latest CA bundle
# hadolint ignore=DL3018
RUN apk add --no-cache ca-certificates


FROM alpine:${ALPINE_VERSION} AS tmp

ARG DATA_PATH
RUN mkdir -m 0750 ${DATA_PATH} && chown 20000:20000 ${DATA_PATH}


## Unix Filesystem Hierarchy Standard 3 (FHS 3.0) directories
## https://refspecs.linuxfoundation.org/FHS_3.0/fhs/index.html
FROM alpine:${ALPINE_VERSION} AS fhs

RUN <<EOF
mkdir -m 0755 \
  /rootfs \
  /rootfs/bin \
  /rootfs/bin/include \
  /rootfs/bin/lib \
  /rootfs/bin/local \
  /rootfs/bin/sbin \
  /rootfs/bin/share/man \
  /rootfs/bin/share/misc \
  /rootfs/dev \
  /rootfs/etc \
  /rootfs/etc/opt \
  /rootfs/lib \
  /rootfs/media \
  /rootfs/mnt \
  /rootfs/opt \
  /rootfs/run \
  /rootfs/sbin \
  /rootfs/srv \
  /rootfs/usr \
  /rootfs/usr/include \
  /rootfs/usr/lib \
  /rootfs/usr/local \
  /rootfs/usr/local/bin \
  /rootfs/usr/local/etc \
  /rootfs/usr/local/games \
  /rootfs/usr/local/include \
  /rootfs/usr/local/lib \
  /rootfs/usr/local/man \
  /rootfs/usr/local/sbin \
  /rootfs/usr/local/share \
  /rootfs/usr/local/src \
  /rootfs/usr/sbin \
  /rootfs/usr/share \
  /rootfs/usr/share/man \
  /rootfs/usr/share/misc \
  /rootfs/var \
  /rootfs/var/cache \
  /rootfs/var/lib \
  /rootfs/var/lib/misc \
  /rootfs/var/lock \
  /rootfs/var/log \
  /rootfs/var/opt \
  /rootfs/var/spool \
  ;

mkdir -m 1777 \
  /rootfs/tmp \
  /rootfs/var/tmp \
  ;

ln -sf /run /var/run

## Linux-specific annex
mkdir -m 0555 \
  /rootfs/proc \
  /rootfs/sys \
  ;
EOF


FROM alpine:${ALPINE_VERSION} AS machine-id

RUN sysctl -n kernel.random.uuid > /etc/machine-id


FROM scratch

WORKDIR /
COPY --from=fhs /rootfs/ /

COPY --from=machine-id /etc/machine-id /etc/machine-id

ARG DATA_PATH
COPY --from=tmp --chown=20000:20000 --chmod=0750 ${DATA_PATH} ${DATA_PATH}

COPY --from=certificates /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY --from=build /src/bin/handler /usr/bin/handler

CMD ["/usr/bin/handler"]

ENV DATA_PATH=${DATA_PATH}
EXPOSE 8080
USER 20000:20000
VOLUME ${DATA_PATH}

LABEL org.opencontainers.image.authors="William Artero <docker@artero.dev>"
LABEL org.opencontainers.image.base.name="scratch"
LABEL org.opencontainers.image.description="anilist custom list provider for sonarr/radarr"
LABEL org.opencontainers.image.documentation="https://github.com/wwmoraes/anilistarr/blob/master/README.md"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/wwmoraes/anilistarr"
LABEL org.opencontainers.image.title="Anilistarr"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/wwmoraes/anilistarr"
LABEL org.opencontainers.image.vendor="William Artero <docker@artero.dev>"
