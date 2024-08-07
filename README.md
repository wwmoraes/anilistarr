# anilistarr

> anilist custom list provider for sonarr/radarr

[![Go Reference](https://pkg.go.dev/badge/github.com/wwmoraes/anilistarr.svg)](https://pkg.go.dev/github.com/wwmoraes/anilistarr)
[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/anilistarr.svg)](https://github.com/wwmoraes/anilistarr/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/anilistarr.svg)](https://github.com/wwmoraes/anilistarr/pulls)
![Codecov](https://img.shields.io/codecov/c/github/wwmoraes/anilistarr)

![GitHub branch status](https://img.shields.io/github/checks-status/wwmoraes/anilistarr/master)
[![Integration](https://github.com/wwmoraes/anilistarr/actions/workflows/integration.yml/badge.svg)](https://github.com/wwmoraes/anilistarr/actions/workflows/integration.yml)
[![Release](https://github.com/wwmoraes/anilistarr/actions/workflows/release.yml/badge.svg)](https://github.com/wwmoraes/anilistarr/actions/workflows/release.yml)
[![Security](https://github.com/wwmoraes/anilistarr/actions/workflows/security.yml/badge.svg)](https://github.com/wwmoraes/anilistarr/actions/workflows/security.yml)
[![Documentation](https://github.com/wwmoraes/anilistarr/actions/workflows/documentation.yml/badge.svg)](https://github.com/wwmoraes/anilistarr/actions/workflows/documentation.yml)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/7718/badge)](https://bestpractices.coreinfrastructure.org/projects/7718)

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/wwmoraes/anilistarr)](https://hub.docker.com/r/wwmoraes/anilistarr)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/wwmoraes/anilistarr?label=image%20version)](https://hub.docker.com/r/wwmoraes/anilistarr)
[![Docker Pulls](https://img.shields.io/docker/pulls/wwmoraes/anilistarr)](https://hub.docker.com/r/wwmoraes/anilistarr)

[![built with nix](https://builtwithnix.org/badge.svg)](https://builtwithnix.org)

---

## 📝 Table of Contents

- [About](#-about)
- [Getting Started](#-getting-started)
- [Deployment](#-deployment)
- [Usage](#-usage)
- [Built Using](#-built-using)
- [TODO](./TODO.md)
- [Contributing](./CONTRIBUTING.md)
- [Authors](#-authors)
- [Acknowledgments](#-acknowledgements)

## 🧐 About

Converts an Anilist user watching list to a custom list format which *arr apps support.

It works by fetching the user info directly from Anilist thanks to its API, and
converts the IDs using community-provided mappings.

Try it out on a live instance at `https://anilistarr.fly.dev/`. For API details
check either the [source Swagger definition](./swagger.yaml) or the generated
[online version here][swagger-ui].

[swagger-ui]: https://editor-next.swagger.io/?url=https%3A%2F%2Fraw.githubusercontent.com%2Fwwmoraes%2Fanilistarr%2Fmaster%2Fswagger.yaml

## 🏁 Getting Started

Clone the repository and use `go run ./cmd/handler/...` to get the REST API up.

## 🔧 Running the tests

Explain how to run the automated tests for this system.

## 🎈 Usage

Configuration in general is a WIP. The code supports distinct storage and cache
options and has built-in support for different caches and stores. The handler
needs flags/configuration file support to allow switching at runtime.

Implemented solutions:

- Cache
  - Badger
  - Bolt (no TTL support tho)
  - Redis
- Store
  - Badger
  - SQL (model generated for SQLite, should work for others but YMMV)

## 🚀 Deployment

The `handler` binary is statically compiled and serves both the REST API and the
telemetry to an OTLP endpoint. Extra requirements depend on which storage and
cache technologies you've chosen; e.g. using SQLite/Bolt requires a database
file. The Docker image provided contains the handler alone, for instance.

## 🔧 Built Using

- [Golang](https://go.dev) - Base language
- [Chi](https://go-chi.io) - net/HTTP-compatible router that doesn't suck
- [genqlient](https://github.com/Khan/genqlient) - type-safe GraphQL client generator
- [xo](https://github.com/xo/xo) - SQL client code generator
- [Open Telemetry](https://opentelemetry.io) - Observability

## 🧑‍💻 Authors

- [@wwmoraes](https://github.com/wwmoraes) - Idea & Initial work

## 🎉 Acknowledgements

- Anilist for their great service and API <https://anilist.gitbook.io/anilist-apiv2-docs/>
- The community for their efforts to map IDs between services
  - <https://github.com/Fribb/anime-lists>
  - <https://github.com/Anime-Lists/anime-lists/>
