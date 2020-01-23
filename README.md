# jsonapi

[![Build Status](https://travis-ci.org/256dpi/jsonapi.svg?branch=master)](https://travis-ci.org/256dpi/jsonapi)
[![Coverage Status](https://coveralls.io/repos/github/256dpi/jsonapi/badge.svg?branch=master)](https://coveralls.io/github/256dpi/jsonapi?branch=master)
[![GoDoc](https://godoc.org/github.com/256dpi/jsonapi?status.svg)](http://godoc.org/github.com/256dpi/jsonapi)
[![Release](https://img.shields.io/github/release/256dpi/jsonapi.svg)](https://github.com/256dpi/jsonapi/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/256dpi/jsonapi)](http://goreportcard.com/report/256dpi/jsonapi)

**A fundamental and extendable JSON API library for Go.**

Package [`jsonapi`](http://godoc.org/github.com/256dpi/jsonapi) provides structures and functions to implement [JSON API](http://jsonapi.org) compatible APIs. The library can be used with any framework and is built on top of the standard Go http library.

## Extensions

### Custom Actions

The package supports the non-standard but widely adopted "custom actions" extension to support the following patterns:

```
GET /posts/highlighted
DELETE /posts/cache
POST /posts/1/publish
DELETE /posts/1/history
```

## Examples

The testing [server](https://github.com/256dpi/jsonapi/blob/master/server.go) implements a basic API server using the standard HTTP package.

## Installation

Get the package using the go tool:

```bash
$ go get -u github.com/256dpi/jsonapi/v2
```

## License

The MIT License (MIT)

Copyright (c) 2016 Joël Gähwiler
