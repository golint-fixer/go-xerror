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

To create a new error, use the `xerror.New` function:

```go
// Defining each error type as a constant string is a good practice.
const ErrorInvalidValueForField = "invalid value for field %v"

// The return type of this function could alternatively be xerror.Error.
// Since it's public and possibly part of a library, we return error.
func ValidateRequest(r *Request) error {
  if r.UserID == "" {
      // "userId" replaces the %v placeholder in the error type string
      // extra arguments such as `request` are instead only attached to the debug objects slice
      return xerror.New(ErrorInvalidValueForField, "userId", request)
  }
  return nil
}
```

Calling the Error interface methods on the newly created error would return the following:

```
err.Error() -> "invalid value for field userId"
err.Debug() -> []interface{}{"userId", request}
err.Stack() -> a slice of strings representing the stack at New call
```

Please note that all arguments besides the format string are appended to the debug objects slice. The library counts how many placeholders are present in the format string and limits the numbers of arguments passed to `fmt.Sprintf` when generating the formatted version.

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

```
err.Error() -> "bad request for URL <contents of r.URL>: unexpected end of file"
err.Debug() -> []interface{}{r.URL, r}
err.Stack() -> a slice of strings representing the stack at Wrap call
```

Calling the Error interface methods on the second error would return the following:

```
err.Error() -> "bad request for URL <contents of r.URL>: malformed request body: invalid character 'b'"
err.Debug() -> []interface{}{r.URL, r, buf
err.Stack() -> a slice of strings representing the stack at the first Wrap call
}
```
