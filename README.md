# XERROR

Package `xerror` extends the functionality of Go's built-in `error` interface: it allows to generate nicely formatted error messages while making it easy to programmatically check for error types, and allowing to propagate additional information such as stack traces and debug values.

Features:

- a stack trace is attached to errors at creation
- additional debug values can be attached to errors for deferred out-of-band logging and reporting
- nice interface for wrapping errors and propagating them at the right level of abstraction
- compatible with Go's `error` interface
- easy to check for error types while generating nicely formatted messages, which include specifics

This package is particularly useful in applications such as API servers, where errors returned to users might contain less detail than those stored in logs. Additionally, the library interoperates well with Go's built-in errors: it allows to easily wrap them for propagation, and it generates errors that implement the `error` interface, making it suitable for use in libraries where the clients do not depend on this package.

It is currently recommended to use the stable v1 branch, which is in maintenance mode. Refer to the `godoc` pages for usage details. More information on branches and future development follows.

## Stable Branch (v1)

[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=v1)](https://travis-ci.org/ibrt/go-xerror?branch=v1)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=v1)](https://coveralls.io/github/ibrt/go-xerror?branch=v1)
[![GoDoc](https://godoc.org/gopkg.in/ibrt/go-xerror.v1/xerror?status.svg)](https://godoc.org/gopkg.in/ibrt/go-xerror.v1/xerror)

- https://github.com/ibrt/go-xerror/tree/v1
- `go get gopkg.in/ibrt/go-xerror.v1/xerror`

## Development Branch (master):

[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=master)](https://travis-ci.org/ibrt/go-xerror?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=master)](https://coveralls.io/github/ibrt/go-xerror?branch=master)
[![GoDoc](https://godoc.org/github.com/ibrt/go-xerror/xerror?status.svg)](https://godoc.org/github.com/ibrt/go-xerror/xerror)

- https://github.com/ibrt/go-xerror
- `go get github.com/ibrt/go-xerror/xerror`
