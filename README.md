# stride
A code generator for OpenAPI v3.x BETA.

[![Documentation][godoc-img]][godoc-url]
[![License][license-img]][license-url]
[![Build Status][action-img]][action-url]
[![Coverage][codecov-img]][codecov-url]
[![Go Report Card][report-img]][report-url]

## Overview

Stride is a command line tool for rapid application development with OpenAPI v3
specification.

## Motivation

There are a lot of tools around the OpenAPI initiative. However neither of them
gives comprehensive solution for rapid application development. The goal is
to provide a tool that you can use to generate production ready server
applications.

Doesn't matter whether you are starting from just concept, prototype or
building your v1 product. The tool provides functionality for editing and
viewing of OpenAPI specification. It generates a production ready scaffold
projects. Also it can start a mock server. So you can quickly start building
your front-end app without waiting for your actual back-end.

Thanks to the contributors of the following projects that empower `stride`:
- [kin-openapi](https://github.com/getkin/kin-openapi)
- [swagger-editor](https://github.com/swagger-api/swagger-editor)
- [swagger-ui](https://github.com/swagger-api/swagger-ui)

## Introduction

It provides the following subcommands:

```bash
NAME:
   stride - OpenAPI viewer, editor, generator, validator and mocker
USAGE:
   stride [global options]
COMMANDS:
     edit      Edit an OpenAPI specification in the browser
     view      Shows an OpenAPI specification in the browser
     mock      Runs a mock server from an OpenAPI specification
     generate  Generates a project from an OpenAPI specification
     validate  Validates an OpenAPI specification
     help, h   Shows a list of commands or help for one command

OPTIONS:
   --version, -v  prints the version
   --help, -h     shows help
```

Right now `stride` supports only `golang`. There are a few limitations that the
generator does not support for now. The following features are not supported:

- Inheritance and Polymorphism
- OneOf, AnyOf, AllOf and Not (there are some limitations due to the language constraints)
- Links
- Callbacks
- Authentication

## Road map

- [x] Golang generator (in testing phase)
- [ ] Enable a mock server for given OpenAPI specification
- [x] Download the OpenAPI specification from different sources (local, s3, git and etc.)
- [x] Support for Dictionaries, Hash Maps and Associative Arrays in Golang
- [x] Support for `application/xml` and `application/x-www-form-urlencoded`
- [ ] Improve the OpenAPI validation reports
- [ ] Allow implementation of 3rd party generators in other languages (via GRPC)

## Installation

There are a few ways to install it:

#### GitHub

```console
$ go get -u github.com/phogolabs/stride
$ go install github.com/phogolabs/stride/cmd/stride
```

#### Homebrew (for Mac OS X)

```console
$ brew tap phogolabs/tap
$ brew install stride
```

## Contributing

We are open for any contributions. Just fork the
[project](https://github.com/phogolabs/stride).

[report-img]: https://goreportcard.com/badge/github.com/phogolabs/stride
[report-url]: https://goreportcard.com/report/github.com/phogolabs/stride
[logo-author-url]: https://www.freepik.com/free-vector/abstract-cross-logo-template_1185919.htm
[logo-license]: http://creativecommons.org/licenses/by/3.0/
[codecov-url]: https://codecov.io/gh/phogolabs/stride
[codecov-img]: https://codecov.io/gh/phogolabs/stride/branch/master/graph/badge.svg
[action-img]: https://github.com/phogolabs/stride/workflows/pipeline/badge.svg
[action-url]: https://github.com/phogolabs/stride/actions
[godoc-url]: https://godoc.org/github.com/phogolabs/stride
[godoc-img]: https://godoc.org/github.com/phogolabs/stride?status.svg
[license-img]: https://img.shields.io/badge/license-MIT-blue.svg
[license-url]: LICENSE
