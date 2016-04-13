# XERROR v2 (stable branch)

[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=v2)](https://travis-ci.org/ibrt/go-xerror?branch=v2)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=v2)](https://coveralls.io/github/ibrt/go-xerror?branch=v2)
[![GoDoc](https://godoc.org/gopkg.in/ibrt/go-xerror.v2/xerror?status.svg)](https://godoc.org/gopkg.in/ibrt/go-xerror.v2/xerror)

```
go get gopkg.in/ibrt/go-xerror.v2/xerror
```

### Overview

Package `xerror` extends the functionality of Go's built-in `error` interface: it allows to generate nicely formatted error messages while making it easy to programmatically check for error types, and allowing to propagate additional information such as stack traces and debug values.

This package is particularly useful in applications such as API servers, where errors returned to users might contain less detail than those stored in logs. Additionally, the library interoperates well with Go's built-in errors: it allows to easily wrap them for propagation, and it generates errors that implement the `error` interface, making it suitable for use in libraries where the clients do not depend on this package.

##### Features

- a stack trace is attached to errors at creation
- additional debug values can be attached to errors for deferred out-of-band logging and reporting
- nice interface for wrapping errors and propagating them at the right level of abstraction
- compatible with Go's `error` interface
- easy to check for error types while generating nicely formatted messages, which include specifics

### How To

We will now learn how to create errors, propagate them, check for error types, access stack traces and debug objects, and interoperate with the standard Go library. This how-to attempts to describe and clarify the best practices for error handling in Go.

##### Creating a new error

- https://godoc.org/gopkg.in/ibrt/go-xerror.v2/xerror#New

```go
// Defining each error type as a constant string is a good practice.
const ErrorInvalidValueForField = "invalid value for field %v"

// The return type of this function could alternatively be xerror.Error.
// Since it's public and possibly part of a library, we return error.
func ValidateRequest(r *Request) error {
  if r.UserID == "" {
      // "userId" replaces the %v placeholder in the error string
      // extra arguments such as `request` are instead only attached to the debug objects slice
      return xerror.New(ErrorInvalidValueForField, "userId", request)
  }
  return nil
}
```
