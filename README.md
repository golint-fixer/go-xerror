# go-xerror
[![Build Status](https://api.travis-ci.org/ibrt/go-xerror.svg?branch=master)](https://travis-ci.org/ibrt/go-xerror?branch=master)
[![Coverage Status](https://coveralls.io/repos/github/ibrt/go-xerror/badge.svg?branch=master)](https://coveralls.io/github/ibrt/go-xerror?branch=master)

Package xerror extends the functionality of Go's built-in error interface, in several ways:

- errors carry a stack trace of the location where they were first created
- it is possible to attach debug objects to errors for later reporting
- errors hold a list of messages, and can be wrapped by prepending new messages to the list
- errors can natively be treated as Go built-in errors and serialized to a string message or JSON representation

### Installation and Documentation

Latest stable version (v1):

[![GoDoc](https://godoc.org/gopkg.in/ibrt/go-xerror.v1/xerror?status.svg)](https://godoc.org/gopkg.in/ibrt/go-xerror.v1/xerror)

```go get gopkg.in/ibrt/go-xerror.v1/xerror```

Development version (might break at any time):

[![GoDoc](https://godoc.org/github.com/ibrt/go-xerror/xerror?status.svg)](https://godoc.org/github.com/ibrt/go-xerror/xerror)

```go get github.com/ibrt/go-xerror/xerror```

### Instructions

To create a new error, use the New function:

```go
err := xerror.New("first message", "second message")
```

When this error is converted to string using the Error method from Go's error interface, the following is returned:

```go
"first message: second message"
```
  
To create an augmented error given a Go error, use the Wrap function:

```go
if _, err := pkg.Method(...); err != nil {
	return xerror.Wrap(err).WithMessages("unable to execute Method")
}
```

If the given error is actually of type Error, the Error is immediately returned unmodified.

Errors are immutable, but modified copies of them can be obtained using WithMessages and WithDebug:

```go
return xerror.Wrap(err).WithMessages("unable to perform action").WithDebug(ctx, req)
```

Error also provides methods for determining its type, by matching all messages or just the outermost one:

```go
err := xerror.New("m2", "m1")
err.Is("m2") // true
err.Is("m1") // false
err.IsPattern(regexp.MustCompile("2$")) // true
err.IsPattern(regexp.MustCompile("1$")) // false
err.Contains("m2") // true
err.Contains("m1") // true
err.ContainsPattern(regexp.MustCompile("2$") // true
err.ContainsPattern(regexp.MustCompile("1$") // true
```
