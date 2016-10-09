# XERROR v2 (stable branch)

[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=v2)](https://travis-ci.org/ibrt/go-xerror?branch=v2)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=v2)](https://coveralls.io/github/ibrt/go-xerror?branch=v2)
[![Go Report Card](https://goreportcard.com/badge/gopkg.in/ibrt/go-xerror.v2)](https://goreportcard.com/report/gopkg.in/ibrt/go-xerror.v2)
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
- easy to check for error types while generating nicely formatted messages
- compatible with Go's `error` interface

### How To

We will now learn how to create errors, propagate them, check for error types, access stack traces and debug objects, and interoperate with the standard Go library. This how-to also attempts to describe and clarify some best practices for error handling in Go.

##### Creating a new error

To create a new error, use the `xerror.New` function:

```go
const ErrorInvalidValueForField = "invalid value for field %v"

func ValidateRequest(r *Request) error {
  if r.UserID <= 0 {
      return xerror.New(ErrorInvalidValueForField, "userId", request)
  }
  return nil
}
```

All arguments of `New` besides the format string are appended to the debug objects slice. The library counts how many placeholders are present in the format string and limits the numbers of arguments passed to `fmt.Sprintf` when generating the formatted version. Calling the Error interface methods on the newly created error would return the following:

```go
err.Error() // -> "invalid value for field userId"
err.Debug() // -> []interface{}{"userId", request}
err.Stack() // -> a slice of strings representing the stack when New is called
```

##### Propagating errors

Errors are usually propagated up the call stack as return values. It is often desirable to wrap them with information at the right level of abstraction, but the standard Go library doesn't provide a good way to do so. The `xerror.Wrap` function can be used for this purpose, as illustrated below:

```go
const (
  ErrorMalformedRequestBody = "malformed request body"
  ErrorBadRequest = "bad request for URL %v"
)

func ParseRequest(buf []byte) (*Request, error) {
  req := &Request{}
  if err := json.Unmarshal(buf, req); err != nil {
    return nil, xerror.Wrap(err, ErrorMalformedRequestBody, buf)
  }
  return req, nil
}

func HandleRequest(r *http.Request) (*Response, error) {
  buf, err := ioutil.ReadAll(r.body)
  defer r.Body.Close()
  if err != nil {
    return nil, xerror.Wrap(err, ErrorBadRequest, r.URL, r) // first error
  }
  req, err := ParseRequest(buf)
  if err != nil {
    return nil, xerror.Wrap(err, ErrorBadRequest, r.URL, r) // second error
  }
  
  ...
}
```

Calling the Error interface methods on the first error would return the following:

```go
err.Error() // -> "bad request for URL http://some-url: unexpected end of file"
err.Debug() // -> []interface{}{r.URL, r}
err.Stack() // -> a slice of strings representing the stack at Wrap call
```

Calling the Error interface methods on the second error would return the following:

```go
err.Error() // -> "bad request for URL http://some-url: malformed request body: invalid character 'b'"
err.Debug() // -> []interface{}{r.URL, r, buf}
err.Stack() // -> a slice of strings representing the stack at the first Wrap call
```

##### Determining the type of an error

This library provides functions for determining error types: `Is` and `Contains`. They exist both as top-level package functions and as methods on the `Error` interface. Error type checking in Go is usually done by storing error messages a string constants, and performing string comparisons. Unfortunately this technique doesn't work well when used together with `fmt.Errorf`, as the generated error string is not equal to the original format string. These functions instead perform the comparison on the format string, allowing to generate clearer error messages while retaining the ability to check for error types.

Let's consider the second error from the _Propagating errors_ section. Here is the result of some sample calls:

```go
err.Is(ErrorBadRequest) // -> true
err.Is(ErrorMalformedRequestBody) // -> false
err.Contains(ErrorBadRequest) // -> true
err.Contains(ErrorMalformedRequestBody) // -> true
```

In other words, `Is` only compares the format string with the outermost error in the wrap chain, while `Contains` performs the comparison on all wrapped errors. The top-level functions work similarly, but they accept any kind of `error` argument. If the given `error` is actually a `xerror.Error`, they are equivalent to calling the corresponding methods on the interface, otherwise they perform the comparison on the string version of the given error.

```go
xerror.Is(secondError, ErrorBadRequest) // -> true
xerror.Is(secondError, ErrorMalformedRequestBody) // -> false
xerror.Is(errors.New(ErrorBadRequest), ErrorBadRequest) // -> true
xerror.Contains(errors.New(ErrorBadRequest), ErrorBadRequest) // -> true
```

##### Reporting and displaying errors

The `xerror.Error` interface extends `error`, `json.Marshaler`, and `fmt.GoStringer`. It is possible to obtain string representations of errors for various use cases:

- calling `err.Error()` or formatting as `%s` or `%v`returns a short string
- serializing to JSON or formatting as `%#v` returns a long string

This is an example of short string:

```
bad request: malformed request body: invalid character 'b'
```

This is an example of long string (actually on a single line):

```
{
  "message": "bad request: malformed request body: invalid character 'b'",
  "debug": [
    "d2",
    "d1"
  ],
  "stack":[
    "/path/to/file1.go:49 (0x8448b)",
    "/path/to/file2.go:198 (0x8448b)",
    ...
  ]
}
```
