# anilistarr

> anilist custom list provider for sonarr/radarr

[![Go Reference](https://pkg.go.dev/badge/github.com/wwmoraes/anilistarr.svg)](https://pkg.go.dev/github.com/wwmoraes/anilistarr)
[![GitHub Issues](https://img.shields.io/github/issues/wwmoraes/anilistarr.svg)](https://github.com/wwmoraes/anilistarr/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/wwmoraes/anilistarr.svg)](https://github.com/wwmoraes/anilistarr/pulls)

[![Integration](https://github.com/wwmoraes/anilistarr/actions/workflows/integration.yml/badge.svg)](https://github.com/wwmoraes/anilistarr/actions/workflows/integration.yml)
[![pre-commit.ci status](https://results.pre-commit.ci/badge/github/wwmoraes/anilistarr/master.svg)](https://results.pre-commit.ci/latest/github/wwmoraes/anilistarr/master)

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_anilistarr&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=wwmoraes_anilistarr)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_anilistarr&metric=coverage)](https://sonarcloud.io/summary/new_code?id=wwmoraes_anilistarr)
[![Maintainability Rating](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_anilistarr&metric=sqale_rating)](https://sonarcloud.io/summary/new_code?id=wwmoraes_anilistarr)
[![Security Rating](https://sonarcloud.io/api/project_badges/measure?project=wwmoraes_anilistarr&metric=security_rating)](https://sonarcloud.io/summary/new_code?id=wwmoraes_anilistarr)

[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fwwmoraes%2Fanilistarr.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fwwmoraes%2Fanilistarr?ref=badge_shield)
[![CII Best Practices](https://bestpractices.coreinfrastructure.org/projects/7718/badge)](https://bestpractices.coreinfrastructure.org/projects/7718)

[![Docker Image Size (latest semver)](https://img.shields.io/docker/image-size/wwmoraes/anilistarr)](https://hub.docker.com/r/wwmoraes/anilistarr)
[![Docker Image Version (latest semver)](https://img.shields.io/docker/v/wwmoraes/anilistarr?label=image%20version)](https://hub.docker.com/r/wwmoraes/anilistarr)
[![Docker Pulls](https://img.shields.io/docker/pulls/wwmoraes/anilistarr)](https://hub.docker.com/r/wwmoraes/anilistarr)

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

## 🏁 Getting Started

Clone the repository and use `go run ./cmd/handler/...` to get the REST API up.

## 🔧 Running the tests

Explain how to run the automated tests for this system.

## 🎈 Usage

Configuration in general is a WIP. The code supports distinct storage and cache
options and even has built-in support for Redis and Bolt as caches already.
The handler needs flags/configuration file support to allow switching at
runtime.

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
