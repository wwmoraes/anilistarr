ARG GOLANG_VERSION=1.20
FROM golang:${GOLANG_VERSION}-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .
ARG VERSION
RUN go generate ./... && go build -o ./bin/handler ./cmd/handler/...


FROM scratch

LABEL org.opencontainers.image.authors="William Artero <docker@artero.dev>"
LABEL org.opencontainers.image.description="anilist custom list provider for sonarr/radarr"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="Anilistarr"
LABEL org.opencontainers.image.vendor="William Artero <docker@artero.dev>"

CMD ["/usr/local/bin/handler"]
EXPOSE 8080

ARG GOLANG_VERSION
LABEL org.opencontainers.image.base.name="golang:${GOLANG_VERSION}-alpine"
ARG VERSION
LABEL org.opencontainers.image.version="${VERSION}"

COPY --from=build /src/bin/handler /usr/local/bin/handler

USER 20000:20000
